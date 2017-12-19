package ghost

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/lunixbochs/ghostrace/ghost/memio"
	"github.com/lunixbochs/ghostrace/ghost/process"
	"github.com/lunixbochs/ghostrace/ghost/sys"
	"github.com/lunixbochs/ghostrace/ghost/sys/call"
)

type execCb func(e *Event) (bool, bool)

func isStopSig(sig syscall.Signal) bool {
	return sig == syscall.SIGSTOP || sig == syscall.SIGTSTP || sig == syscall.SIGTTIN || sig == syscall.SIGTTOU
}

type Tracer interface {
	ExecFilter(cb execCb)
	Spawn(cmd string, args ...string) (chan *Event, error)
	Trace(pid int) (chan *Event, error)
}

type LinuxTracer struct {
	execFilter execCb
}

func NewTracer() Tracer {
	return &LinuxTracer{}
}

func (t *LinuxTracer) ExecFilter(cb execCb) {
	t.execFilter = cb
}

func (t *LinuxTracer) Spawn(cmd string, args ...string) (chan *Event, error) {
	pid, err := syscall.ForkExec(cmd, args, &syscall.ProcAttr{
		Sys:   &syscall.SysProcAttr{Ptrace: true},
		Files: []uintptr{0, 1, 2},
	})
	if err != nil {
		return nil, err
	}
	return t.traceProcess(pid, true)
}

func (t *LinuxTracer) Trace(pid int) (chan *Event, error) {
	return t.traceProcess(pid, false)
}

func (t *LinuxTracer) traceProcess(pid int, spawned bool) (chan *Event, error) {
	ret := make(chan *Event)
	errChan := make(chan error)
	go func() {
		// rate limit the sigstop loop hack
		stopToken := make(chan int, 10)
		go func() {
			for {
				stopToken <- 1
				time.Sleep(50 * time.Millisecond)
			}
		}()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(ret)

		spawnChild := -1
		if spawned {
			spawnChild = pid
		} else {
			if err := syscall.PtraceAttach(pid); err != nil {
				errChan <- err
				return
			}
		}
		topPid := pid
		errChan <- nil

		first := true
		table := make(map[int]*tracedProc)
		// we need to catch interrupts so we don't leave other processes in a bad state
		// TODO: make the interrupt catching behavior optional (but default)?
		var interrupted syscall.Signal
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGQUIT)
		go func() {
			for sig := range signalChan {
				// TODO: send an interrupt event back over the channel?
				// otherwise just make the other side also listen for interrupts
				interrupted = sig.(syscall.Signal)
				// interrupt the main loop's Wait4
				syscall.Kill(topPid, syscall.SIGSTOP)
			}
		}()
		var cleanup = func(exitSig syscall.Signal) {
			if spawnChild >= 0 {
				syscall.Kill(pid, syscall.SIGCONT)
				syscall.Kill(pid, exitSig)
			}
			if interrupted != 0 {
				for _, traced := range table {
					traced.Detach()
				}
			}
		}
		defer func() {
			if recover() != nil {
				interrupted = syscall.SIGSEGV
				cleanup(syscall.SIGTERM)
			}
		}()
		for interrupted == 0 {
			var status syscall.WaitStatus
			pid, err := syscall.Wait4(-1, &status, syscall.WALL, nil)
			if err != nil {
				if err == syscall.EINTR {
					continue
				} else if err == syscall.ECHILD && len(table) == 0 {
					break
				}
				fmt.Println("DEBUG:", err)
				break
			}
			sig := status.StopSignal()
			traced, ok := table[pid]
			if ok {
				if !traced.EatOneSigstop && isStopSig(sig) {
					// TODO: need to be smarter about rate limiting this
					<-stopToken
					traced.StopSig = sig
					if err = syscall.PtraceSyscall(pid, int(sig)); err != nil && err != syscall.ESRCH {
						break
					}
					continue
				} else if traced.StopSig != 0 {
					if sig == syscall.SIGCONT {
						traced.StopSig = 0
					} else if sig != 0 {
						if err = syscall.PtraceSyscall(pid, int(traced.StopSig)); err != nil && err != syscall.ESRCH {
							break
						}
						continue
					}
				}
				if status.Exited() {
					// process exit
					ret <- &Event{
						Process: traced.Process,
						Exit:    true,
					}
					delete(table, pid)
					continue
				}
			} else {
				// set up new pid
				proc, err := process.FindPid(pid)
				if err != nil {
					fmt.Println("DEBUG:", err)
					continue
				}
				t, err := newTracedProc(proc, !first || !spawned)
				if err != nil {
					fmt.Println("DEBUG:", err)
					continue
				}
				first = false
				traced = t
				table[pid] = t
			}
			if status.TrapCause() != -1 {
				// handle PTRACE_EVENT_*
				sig = syscall.Signal(0)
			} else if sig == syscall.SIGSTOP && traced.EatOneSigstop {
				// handle first SIGSTOP
				traced.EatOneSigstop = false
				sig = syscall.Signal(0)
			} else if sig != syscall.SIGTRAP|0x80 && interrupted != 0 {
				// we can interrupt if it's NOT a syscall delivery
				break
			} else if sig == syscall.SIGTRAP|0x80 {
				// handle a syscall delivery
				var sc sys.Syscall
				var err error
				if sc, err = traced.Syscall(); err != nil {
					fmt.Println("DEBUG:", err)
					continue
				}
				if sc != nil {
					ret <- &Event{
						Process: traced.Process,
						Syscall: sc,
					}
					// TODO: need to update the proc's exe/cmdline after execve
					// maybe add a proc.Reset()?
					if _, ok := sc.(*call.Execve); ok && t.execFilter != nil {
						keepParent, followChild := t.execFilter(&Event{
							Process: traced.Process,
							Syscall: sc,
						})
						if !keepParent {
							parent := traced.Process.Parent()
							if parent != nil {
								pid := parent.Pid()
								if tracedParent, ok := table[pid]; ok {
									tracedParent.Detach()
									delete(table, pid)
								}
							}
						}
						if !followChild {
							traced.Detach()
							delete(table, pid)
						}
					}
				}
				sig = syscall.Signal(0)
			}
			// TODO: send events upstream for signals

			// continue and pass signal (which might be zero from earlier code)
			if err = syscall.PtraceSyscall(pid, int(sig)); err != nil && err != syscall.ESRCH {
				break
			}
		}
		exitSig := interrupted
		if exitSig == 0 {
			exitSig = syscall.SIGTERM
		}
		cleanup(exitSig)
	}()
	err := <-errChan
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type tracedProc struct {
	Process       process.Process
	Codec         *sys.Codec
	StopSig       syscall.Signal
	NewSyscall    bool
	SavedRegs     syscall.PtraceRegs
	EatOneSigstop bool
}

func newTracedProc(proc process.Process, eat bool) (*tracedProc, error) {
	pid := proc.Pid()
	options := syscall.PTRACE_O_TRACECLONE | syscall.PTRACE_O_TRACEFORK | syscall.PTRACE_O_TRACEVFORK | syscall.PTRACE_O_TRACESYSGOOD
	if err := syscall.PtraceSetOptions(pid, options); err != nil {
		return nil, err
	}
	var readMem = func(p []byte, addr uint64) (int, error) {
		return syscall.PtracePeekData(pid, uintptr(addr), p)
	}
	var writeMem = func(p []byte, addr uint64) (int, error) {
		return syscall.PtracePokeData(pid, uintptr(addr), p)
	}
	codec, err := sys.NewCodec(sys.ARCH_X86_64, sys.OS_LINUX, memio.NewMemIO(readMem, writeMem))
	if err != nil {
		return nil, err
	}
	return &tracedProc{
		Process:       proc,
		Codec:         codec,
		NewSyscall:    true,
		EatOneSigstop: eat,
	}, nil
}

func (t *tracedProc) Syscall() (ret sys.Syscall, err error) {
	pid := t.Process.Pid()
	if t.NewSyscall {
		t.NewSyscall = false
		if err = syscall.PtraceGetRegs(pid, &t.SavedRegs); err != nil {
			return
		}
	} else {
		t.NewSyscall = true
	}
	name := t.Codec.GetName(int(t.SavedRegs.Orig_rax))
	if t.NewSyscall != (name == "execve") {
		regs := &t.SavedRegs
		var newRegs syscall.PtraceRegs
		syscall.PtraceGetRegs(pid, &newRegs)
		args := []uint64{regs.Rdi, regs.Rsi, regs.Rdx, regs.R10, regs.R8, regs.R9}
		sc, err := t.Codec.DecodeRet(int(regs.Orig_rax), args, newRegs.Rax)
		if err != nil {
			fmt.Println(err)
		} else {
			ret = sc
		}
	}
	return
}

func (t *tracedProc) Detach() error {
	pid := t.Process.Pid()
	// if we're expecting a SIGSTOP, skip this and wait for it
	// otherwise the child will be detached into a an immediate SIGSTOP which is awkward for everyone
	if !t.EatOneSigstop {
		if err := syscall.PtraceDetach(pid); err == nil || err != nil && err != syscall.ESRCH {
			return err
		}
		if err := syscall.Kill(pid, 0); err != nil && err != syscall.ESRCH {
			return err
		}
		if err := syscall.Kill(pid, syscall.SIGSTOP); err != nil && err != syscall.ESRCH {
			return err
		}
	}
	for {
		var status syscall.WaitStatus
		if _, err := syscall.Wait4(pid, &status, syscall.WALL, nil); err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		if !status.Stopped() {
			break
		}
		// Linux wants a SIGSTOP before you detach from a process
		sig := status.StopSignal()
		if sig == syscall.SIGSTOP {
			syscall.PtraceDetach(pid)
			break
		}
		if sig&syscall.SIGTRAP != 0 {
			sig = syscall.Signal(0)
		}
		if err := syscall.PtraceCont(pid, int(sig)); err != nil && err != syscall.ESRCH {
			return err
		}
	}
	return nil
}

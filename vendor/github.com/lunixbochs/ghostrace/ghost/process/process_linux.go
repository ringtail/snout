package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var numRe = regexp.MustCompile(`^\d+$`)

type LinuxProcess struct {
	process
	exe     string
	cmdline []string
}

func (p *LinuxProcess) String() string {
	return fmt.Sprintf("<pid: %d cmdline: %s>", p.Pid(), strings.Join(p.Cmdline(), ", "))
}

func (p *LinuxProcess) Exe() string {
	return p.exe
}

func (p *LinuxProcess) Cmdline() []string {
	return p.cmdline
}

func get(pid int) (Process, error) {
	dir, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return nil, err
	}
	stat := dir.Sys().(*syscall.Stat_t)
	exe, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	rawcmdline, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	cmdline := strings.Split(strings.TrimRight(string(rawcmdline), "\x00"), "\x00")
	return &LinuxProcess{
		process: process{
			pid: pid,
			uid: int(stat.Uid),
			gid: int(stat.Gid),
		},
		exe:     exe,
		cmdline: cmdline,
	}, nil
}

func List() (ProcessList, error) {
	proc, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	pl := make([]Process, 0, len(proc))
	for _, p := range proc {
		if p.IsDir() && numRe.Match([]byte(p.Name())) {
			i, _ := strconv.Atoi(p.Name())
			ps, _ := get(i)
			pl = append(pl, ps)
		}
	}
	return pl, nil
}

var statusRe = regexp.MustCompile(`(?im)^ppid:\s*(\d+)$`)

func (p *LinuxProcess) Parent() Process {
	status, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", p.pid))
	if err != nil {
		return nil
	}
	sub := statusRe.FindSubmatch(status)
	if len(sub) >= 2 {
		ppid, _ := strconv.Atoi(string(sub[1]))
		p, _ := get(ppid)
		return p
	}
	return nil
}

func (lp *LinuxProcess) Children() ProcessList {
	list, _ := Filter(func(p Process) bool {
		parent := p.Parent()
		return parent != nil && parent.Pid() == lp.pid
	})
	return list
}

package process

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <sys/sysctl.h>

struct kinfo_proc *process_list(int *count) {
    size_t size;
    int params[] = {CTL_KERN, KERN_PROC, KERN_PROC_ALL};
    if (sysctl(params, 3, NULL, &size, NULL, 0) < 0) {
        return NULL;
    }
    struct kinfo_proc *ret = NULL;
    if (! (ret = malloc(size))) {
        return NULL;
    }
    if (sysctl(params, 3, ret, &size, NULL, 0) < 0) {
        free(ret);
        return NULL;
    }
    *count = size / sizeof(struct kinfo_proc);
    return ret;
}

char *process_args(int pid, int *argc, size_t *size_out) {
    size_t size;
    int params[] = {CTL_KERN, KERN_PROCARGS2, pid};
    // get argc
    if (sysctl(params, 3, NULL, &size, NULL, 0) < 0) {
        return NULL;
    }
    char *buf = malloc(size);
    if (sysctl(params, 3, buf, &size, NULL, 0) < 0) {
        free(buf);
        return NULL;
    }
    *argc = *(int *)buf;
    free(buf);

    params[1] = KERN_PROCARGS;
    // get argv
    if (sysctl(params, 3, NULL, &size, NULL, 0) < 0) {
        return NULL;
    }
    buf = malloc(size);
    if (sysctl(params, 3, buf, &size, NULL, 0) < 0) {
        free(buf);
        return NULL;
    }
    *size_out = size;
    return buf;
}
*/
import "C"

func charToByte(cc []C.char) []byte {
	tmp := make([]byte, len(cc))
	for i, v := range cc {
		tmp[i] = byte(v)
	}
	return bytes.Trim(tmp, "\x00")
}

func cstr(cc []C.char) string {
	return string(bytes.SplitN(charToByte(cc), []byte{0}, 2)[0])
}

type DarwinProcess struct {
	process
	ppid int
	comm string
}

func (p *DarwinProcess) argv() []string {
	var argc C.int
	var size C.size_t
	buf := C.process_args(C.int(p.pid), &argc, &size)
	if buf == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(buf))

	count := int(argc)
	tmp := string(charToByte((*[1 << 30]C.char)(unsafe.Pointer(buf))[:size:size]))
	tmpSplit := strings.SplitN(tmp, "\x00", 2)
	exe, tmp := tmpSplit[0], strings.TrimLeft(tmpSplit[1], "\x00")
	argv := strings.SplitN(tmp, "\x00", count+1)[:count]
	return append([]string{exe}, argv...)
}

func (p *DarwinProcess) String() string {
	return fmt.Sprintf("<pid %d cmdline: %s", p.Pid(), strings.Join(p.Cmdline(), " "))
}

func (p *DarwinProcess) Exe() string {
	argv := p.argv()
	if len(argv) >= 1 {
		return argv[0]
	}
	return p.comm
}

func (p *DarwinProcess) Cmdline() []string {
	argv := p.argv()
	if len(argv) > 1 {
		return argv[1:]
	}
	return nil
}

func get(pid int) (Process, error) {
	return getFallback(pid)
}

func List() (ProcessList, error) {
	var count C.int
	tmp := C.process_list(&count)
	defer C.free(unsafe.Pointer(tmp))
	if tmp == nil {
		return nil, errors.New("Could not retrieve process list.")
	}
	kinfo := (*[1 << 30]C.struct_kinfo_proc)(unsafe.Pointer(tmp))[:count:count]
	pl := make([]Process, 0, count)
	for i := 0; i < int(count); i++ {
		proc := kinfo[i]
		ps := &DarwinProcess{
			process: process{
				pid: int(proc.kp_proc.p_pid),
				uid: int(proc.kp_eproc.e_pcred.p_ruid),
				gid: int(proc.kp_eproc.e_pcred.p_rgid),
			},
			ppid: int(proc.kp_eproc.e_ppid),
			comm: cstr(proc.kp_proc.p_comm[:]),
		}
		pl = append(pl, ps)
	}
	return pl, nil
}

func (p *DarwinProcess) Parent() Process {
	ps, _ := FindPid(p.ppid)
	return ps
}

func (dp *DarwinProcess) Children() ProcessList {
	list, _ := Filter(func(p Process) bool {
		parent := p.Parent()
		return parent != nil && parent.Pid() == dp.pid
	})
	return list
}

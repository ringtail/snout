package process

import (
	"fmt"
	"strings"
)

type process struct {
	pid, uid, gid int
}

func (p *process) Pid() int {
	return p.pid
}

func (p *process) Uid() int {
	return p.uid
}

func (p *process) Gid() int {
	return p.gid
}

type Process interface {
	Pid() int
	Exe() string
	Cmdline() []string
	Uid() int
	Gid() int
	Parent() Process
	Children() ProcessList
	String() string
}

type Match struct {
	Name string
}

func Filter(cb func(Process) bool) (ProcessList, error) {
	list, err := List()
	if err != nil {
		return nil, err
	}
	return list.Filter(cb), nil
}

func FindName(name string) (ProcessList, error) {
	return Filter(func(p Process) bool {
		return strings.Contains(p.Cmdline()[0], name)
	})
}

func FindPid(pid int) (Process, error) {
	return get(pid)
}

func getFallback(pid int) (Process, error) {
	list, err := Filter(func(p Process) bool {
		return p.Pid() == pid
	})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("pid %d not found", pid)
	}
	return list[0], nil
}

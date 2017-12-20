package process

import (
	"fmt"
	"sort"
)

type ProcessList []Process
type byPid ProcessList
type byUid ProcessList

func (a byPid) Len() int           { return len(a) }
func (a byPid) Less(i, j int) bool { return a[i].Pid() < a[j].Pid() }
func (a byPid) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byUid) Len() int           { return len(a) }
func (a byUid) Less(i, j int) bool { return a[i].Uid() < a[j].Uid() }
func (a byUid) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (pl ProcessList) Print(tree bool) {
	sort.Sort(byPid(pl))
	for _, p := range pl {
		fmt.Printf("%+v\n", p)
	}
}

func (pl ProcessList) Filter(cb func(Process) bool) ProcessList {
	out := make(ProcessList, 0, len(pl))
	for _, v := range pl {
		if v != nil && cb(v) {
			out = append(out, v)
		}
	}
	return out
}

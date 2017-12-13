package disk

import (
	"testing"
)
/***
map[Available:14879744000
	Capacity:248391270400
	Usage:233249382400
	Inodes:60642398
	InodesFree:3632750
	InodesUsed:57009648]
 */

func Test_DiskCollector_Name(t *testing.T) {
	dc := &DiskCollector{}
	if dc.Name() == DISK_STATUS {
		t.Log("pass TestDiskCollector_Name")
		return
	}
	t.Errorf("Failed to pass TestDiskCollector_Name because of name is not equals")
}

func Test_DiskCollector_Gather(t *testing.T) {
	dc := &DiskCollector{}
	mc, err := dc.Gather()
	if err != nil {
		t.Errorf("Failed to gather disk info,because of %s", err.Error())
		return
	}
	t.Logf("pass Test_DiskCollector_Gather and metrics is %v", mc.List())
}

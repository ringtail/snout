package disk

import (
	"testing"
)

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

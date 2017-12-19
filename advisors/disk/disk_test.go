package disk

import (
	"testing"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"github.com/ringtail/snout/collectors"
	"github.com/ringtail/snout/collectors/disk"
)

var (
	Metrics_tree *storage.MetricsTree
)

func init() {
	Metrics_tree = &storage.MetricsTree{}
	Metrics_tree.MetricsSection = make(map[string]types.MetricsSection)
	ds, _ := collectors.Cm.Find(disk.DISK_STATUS).Gather()
	Metrics_tree.AddSection(ds)
}

func Test_DiskAdvisor_Name(t *testing.T) {
	da := &DiskAdvisor{}
	if da.Name() == DISK_ADVISOR {
		t.Logf("pass Test_DiskAdvisor_Name")
		return
	}
	t.Error("Failed to pass Test_DiskAdvisor_Name")
}

func Test_DiskAdvisor_Advise(t *testing.T) {
	da := &DiskAdvisor{}
	symptoms := da.Advise()
	if len(symptoms) == 0 {
		t.Skipf("pass Test_DiskAdvisor_Advise because of not enough metrics")
		return
	}
	t.Log("pass Test_DiskAdvisor_Advise")
}

func Test_GetInodeSymptom(t *testing.T) {
	inode_symptom := GetInodeSymptom(Metrics_tree)
	if inode_symptom == nil {
		t.Skipf("pass Test_GetInodeSymptom because of not enough metrics")
		return
	}
	t.Log("pass Test_GetInodeSymptom")
}

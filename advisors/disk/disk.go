package disk

import (
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"github.com/ringtail/snout/advisors"
)

const (
	DISK_ADVISOR = "DISK_ADVISOR"
)

//TODO
type DiskAdvisor struct{}

func (da *DiskAdvisor) Name() string {
	return DISK_ADVISOR
}

func (da *DiskAdvisor) Description() string {
	return ""
}

func (da *DiskAdvisor) Advise() []types.Symptom {
	tree := storage.InternalMetricsTree
	symptoms := make([]types.Symptom, 0)

	if inode_symptom := GetInodeSymptom(tree); inode_symptom != nil {
		symptoms = append(symptoms, inode_symptom)
	}
	return symptoms
}

func init() {
	da := &DiskAdvisor{}
	advisors.Add(da)
}

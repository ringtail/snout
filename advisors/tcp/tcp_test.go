package tcp

import (
	"testing"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"github.com/ringtail/snout/collectors"
	"github.com/ringtail/snout/collectors/system"
	"github.com/ringtail/snout/collectors/netstat"
)

var (
	Metrics_tree *storage.MetricsTree
)

func init() {
	Metrics_tree = &storage.MetricsTree{}
	Metrics_tree.MetricsSection = make(map[string]types.MetricsSection)
	ks, _ := collectors.Cm.Find(system.KERNEL_SETTINGS).Gather()
	Metrics_tree.AddSection(ks)
	nt, _ := collectors.Cm.Find(netstat.NETSTAT_STATUS).Gather()
	Metrics_tree.AddSection(nt)
}

func TestGetCloseWaitSymptom(t *testing.T) {
	symptom := GetCloseWaitSymptom(Metrics_tree)
	if symptom != nil {
		t.Logf("pass GetCloseWaitSymptom: %s %s", symptom.GetName(), symptom.GetDescription())
	}
	t.Skipf("pass GetCloseWaitSymptom, because not enough metrics")
}

func TestGetTimeWaitSymptom(t *testing.T) {
	symptom := GetTimeWaitSymptom(Metrics_tree)
	if symptom != nil {
		t.Logf("pass TestGetTimeWaitSymptom: %s %s", symptom.GetName(), symptom.GetDescription())
	}
	t.Skipf("pass TestGetTimeWaitSymptom, because not enough metrics")
}

func TestGetPortRangeSymptom(t *testing.T) {
	symptom := GetPortRangeSymptom(Metrics_tree)
	if symptom != nil {
		t.Logf("pass TestGetPortRangeSymptom: %s %s", symptom.GetName(), symptom.GetDescription())
	}
	t.Skipf("pass TestGetPortRangeSymptom, because not enough metrics")
}

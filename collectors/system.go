package collectors

import (
	"errors"
	"github.com/ringtail/snout/types"
	"github.com/ringtail/sysctl"
)

const (
	KERNEL_SETTINGS = "KERNEL_SETTINGS"
)

func init() {
	sc := &SystemCollector{}
	Add(sc)
}

type SystemCollector struct{}

func (sc *SystemCollector) Name() string {
	return KERNEL_SETTINGS
}

func (sc *SystemCollector) Description() string {
	return "Gather kernel settings from /proc/sys"
}

func (sc *SystemCollector) Gather() (types.MetricsSection, error) {
	metrics := sysctl.All()
	if metrics == nil {
		return nil, errors.New("Failed to Gather " + KERNEL_SETTINGS + " metrics")
	}
	return &types.DefaultMetricsSection{
		Name:    KERNEL_SETTINGS,
		Metrics: metrics,
	}, nil
}

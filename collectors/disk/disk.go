package disk

import (
	"github.com/ringtail/snout/types"
	df "github.com/ringtail/go-statfs"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/ringtail/snout/collectors"
	"strconv"
)

const (
	DISK_STATUS = "DISK_STATUS"
)

type DiskCollector struct{}

func (dc *DiskCollector) Name() string {
	return DISK_STATUS
}

func (dc *DiskCollector) Description() string {
	return ""
}

func (dc *DiskCollector) Gather() (types.MetricsSection, error) {
	diskInfo, err := df.GetDiskInfo("/")
	if err != nil {
		log.Warnf("Failed to get disk info,because of %s", err.Error())
	}
	metricsMap := make(map[string]string)
	di := structs.New(diskInfo)
	for _, f := range di.Fields() {
		metricsMap[f.Name()] = strconv.FormatInt(f.Value().(int64), 10)
	}
	return &types.DefaultMetricsSection{
		Name:    DISK_STATUS,
		Metrics: metricsMap,
	}, nil
}

func init() {
	dc := &DiskCollector{}
	collectors.Add(dc)
}

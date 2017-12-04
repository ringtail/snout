package memory

import (
	"github.com/ringtail/snout/types"
	"github.com/ringtail/snout/collectors"
	"os"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func init() {
	mc := &MemoryCollector{}
	collectors.Add(mc)
}

const (
	MEMORY_PROC_STATUS_PATH = "/proc/meminfo"
)

const (
	MEMORY_STATUS = "MEMORY_STATUS"
)

type MemoryCollector struct{}

func (mc *MemoryCollector) Name() string {
	return MEMORY_STATUS
}

func (mc *MemoryCollector) Description() string {
	return ""
}
func (mc *MemoryCollector) Gather() (types.MetricsSection, error) {
	if _, err := os.Stat(MEMORY_PROC_STATUS_PATH); os.IsNotExist(err) {
		log.Errorf("Failed to get meminfo, Because of %v", err.Error())
		return nil, err
	}
	bytes, err := ioutil.ReadFile(MEMORY_PROC_STATUS_PATH)
	if err != nil {
		log.Errorf("Failed to get meminfo from file, Because of %s", err.Error())
		return nil, err
	}
	meminfo := ReadBytesToMap(bytes)
	return &types.DefaultMetricsSection{
		Name:    MEMORY_STATUS,
		Metrics: meminfo,
	}, nil

}

func ReadBytesToMap(bytes []byte) map[string]string {
	meminfo := make(map[string]string)
	linesStr := string(bytes)
	lineArr := strings.Split(linesStr, "\n")
	for _, line := range lineArr {
		if !strings.Contains(line, ":") {
			continue
		}
		line_key_value := strings.Split(line, ":")
		meminfo[line_key_value[0]] = line_key_value[1]
	}
	return meminfo
}

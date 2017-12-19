package cpu

import (
	"github.com/ringtail/snout/collectors"
	log "github.com/Sirupsen/logrus"
	"github.com/ringtail/snout/types"
	"strconv"
	"os/exec"
	"bytes"
	"strings"
	"fmt"
)

const (
	CPU_USAGE_SETTING = "CPU_USAGE"

	CPU_USAGE_THRESHOLD = 70
)

func init() {
	sc := &CPUCollector{}
	collectors.Add(sc)
}

type CPUCollector struct{}

func (sc *CPUCollector) Name() string {
	return CPU_USAGE_SETTING
}

func (sc *CPUCollector) Description() string {
	return "Gather cpu usage from ps aux"
}

type CPUMetricsSection struct {

	Name    string
	Metrics map[string]string
}

func (sc *CPUCollector) Gather() (types.MetricsSection, error) {
	cpuUsage := GatherAllProcessCpuUsage()
	return &types.DefaultMetricsSection{
		Name:    CPU_USAGE_SETTING,
		Metrics: convertCPUToStringMap(cpuUsage),
	}, nil
}



func GatherAllProcessCpuUsage() map[string] float64{
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	processCPUUsages := make(map[string]float64)
	// Skip first line
	out.ReadString('\n')
	for {
		line, err := out.ReadString('\n')
		if err!=nil {
			break;
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range(tokens) {
			if t!="" && t!="\t" {
				ft = append(ft, t)
			}
		}
		pidStr := ft[1]
		_, err = strconv.Atoi(pidStr)
		if err != nil {
			log.Errorf("Failed to parse pid from %s", ft[1])
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Errorf("Failed to parse pid cpuusage from %s", ft[2])
		}
		processCPUUsages[pidStr] = cpu

	}
	return processCPUUsages
}

func convertCPUToStringMap(processCPU map[string] float64) map[string] string {
	processCPUString := make(map[string]string)
	for pid, cpu := range processCPU {
		if cpu > CPU_USAGE_THRESHOLD {
			processCPUString[pid] = fmt.Sprintf("%s", cpu)
		}
	}
	return processCPUString
}


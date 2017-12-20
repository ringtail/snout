package cpu

import (
	"github.com/lunixbochs/ghostrace/ghost"
	"fmt"
	"log"
	"github.com/ringtail/snout/types"
	"time"
	"strconv"
	"strings"
	"github.com/ringtail/snout/collectors"
	"sync"
	"encoding/json"
)

const (

	TRACE_TIME = 5 * time.Second

	MAX_CPU_PROCESS = "MAX_CPU_PROCESS"
)

func init() {
	sc := &CPUSystemCallCollector{}
	collectors.Add(sc)
}

type CPUSystemCallCollector struct{}

func (sc *CPUSystemCallCollector) Name() string {
	return MAX_CPU_PROCESS
}

func (sc *CPUSystemCallCollector) Description() string {
	return "Gather max cpu usage process detail"
}

func (sc *CPUSystemCallCollector) Gather() (types.MetricsSection, error) {
	metric := GatherMaxCPUProcessMetric()
	return &types.DefaultMetricsSection{
		Name:    CPU_USAGE_SETTING,
		Metrics: metric,
	}, nil
}

func getMaxCpuProcess() (int, float64){
	maxCpu := 0.0
	maxCpuPid := ""
	cpuProcess := GatherAllProcessCpuUsage()
	for pid, cpu := range cpuProcess {
		if cpu > maxCpu {
			maxCpu = cpu
			maxCpuPid = pid
		}
	}
	pid, err := strconv.Atoi(maxCpuPid)
	if err != nil {
		log.Fatal(err)
		return 0, 0
	}
	return pid, maxCpu
}

func GatherMaxCPUProcessMetric() map[string]string {
	systemCallMetric := make(map[string]string)
	pid, cpu := getMaxCpuProcess()

	systemCallMetric["PID"] = fmt.Sprintf("%d", pid)
	systemCallMetric["CURRENT_CPU_USAGE"] = fmt.Sprintf("%f", cpu)

	systemCallRecord, err := GatherProcessSystemCallWithinTime(pid, TRACE_TIME)
	if err != nil {
		return systemCallMetric
	}

	result, _ := json.Marshal(systemCallRecord)

	systemCallMetric["SYSTEM_CALL"] = string(result)
	systemCallMetric["SYSTEM_CALL_RECORD_TIME"] = "5s"

	return systemCallMetric
}

func getSystemCallName(systemCall string) string{
	splitTemp := strings.Split(systemCall, "(")
	return splitTemp[0]
}

func GatherProcessSystemCallWithinTime(pid int, maxTime time.Duration) (map[string]int, error) {
	systemCallRecord := make(map[string]int)
	var wg sync.WaitGroup
	traceExit := false

	wg.Add(1)

	// Trace system call of the process
	go func() {
		defer wg.Done()
		tracer := ghost.NewTracer()
		trace, err := tracer.Trace(pid)
		if err != nil {
			log.Fatal(err)
			return
		}
		for sc := range trace {
			if traceExit {
				break
			}
			if sc.Exit {
				break
			}
			result := fmt.Sprintf("%s", sc.Syscall)
			systemCallName := getSystemCallName(result)
			value, ok := systemCallRecord[systemCallName]; if ok {
				systemCallRecord[systemCallName] = value + 1
			}else {
				systemCallRecord[systemCallName] = 1
			}
		}
	}()

	// Limit trace time
	go func() {
		defer wg.Done()
		startTime := time.Now()
		for {
			time.Sleep(1000 * time.Millisecond)
			if time.Now().Sub(startTime) > maxTime {
				traceExit = true
				break
			}
		}
	}()
	wg.Wait()

	return systemCallRecord, nil
}

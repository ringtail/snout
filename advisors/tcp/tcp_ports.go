package tcp

import (
	"fmt"
	"github.com/ringtail/snout/collectors/netstat"
	"github.com/ringtail/snout/collectors/system"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"strconv"
	"strings"
)

func GetPortRangeSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	kernel_settings := metrics_tree.FindSection(system.KERNEL_SETTINGS)
	netstat_status := metrics_tree.FindSection(netstat.NETSTAT_STATUS)
	ports_total_range := kernel_settings.Find("net.ipv4.ip_local_port_range")
	if ports_total_range != "" {
		n := 1
		ports_total_range_slim := strings.TrimFunc(ports_total_range, func(r rune) bool {
			if r == ' ' {
				if n == 1 {
					n = n + 1
					return false
				}
				return true
			}
			return false
		})

		ports_total_arr := strings.Split(ports_total_range_slim, " ")
		max, _ := strconv.Atoi(ports_total_arr[1])
		min, _ := strconv.Atoi(ports_total_arr[0])

		ports_total := max - min
		ports_usage, _ := strconv.Atoi(netstat_status.Find("PORTS_USAGE"))
		if float32(ports_usage) > 0.8*float32(ports_total) {
			desc := fmt.Sprintf("Current system ports range is between %s, but ports total usage is %v", ports_total_range, ports_usage)
			adviseDescs := []string{
				"Ports Usage too much means connections is too much,Please check connection status is in normal status  " +
					"by `netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}' ` ",

				"Please check weather your application has too much 504 or 502 timeout in logs, " +
					"You can also increase port range by `sudo sysctl -w net.ipv4.ip_local_port_range=\"min_port_num max_port_num\"`",
			}
			ports_usage_symptom := types.CreateTextDefaultSymptom(PORTS_USAGE_TOO_MUCH_SYMPTOM, desc, adviseDescs)
			return ports_usage_symptom
		}
	}
	return nil
}
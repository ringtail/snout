package tcp

import (
	"fmt"
	"github.com/ringtail/snout/advisors"
	"github.com/ringtail/snout/collectors/netstat"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"strconv"
	"strings"
	"github.com/ringtail/snout/collectors/system"
)

func init() {
	ta := &TcpAdvisor{}
	advisors.Add(ta)
}

var (
	TCP_ADVISOR                  = "TCP_ADVISOR"
	TIME_WAIT_TOO_MUCH_SYMPTOM   = "TIME_WAIT_TOO_MUCH"
	CLOSE_WAIT_TOO_MUCH_SYMPTOM  = "CLOSE_WAIT_TOO_MUCH"
	PORTS_USAGE_TOO_MUCH_SYMPTOM = "PORTS_USAGE_TOO_MUCH"
)

const (
	MAX_TIME_OUT_CONNECTION   = 100
	MAX_CLOSE_WAIT_CONNECTION = 100
)

type TcpAdvisor struct {
}

func (ta *TcpAdvisor) Name() string {
	return TCP_ADVISOR
}

func (ta *TcpAdvisor) Description() string {
	return ""
}

func (ta *TcpAdvisor) Advise() []types.Symptom {
	tcp_connection_symptoms := handle_tcp_connection()
	return tcp_connection_symptoms
}

func handle_tcp_connection() []types.Symptom {
	tree := storage.InternalMetricsTree
	symptoms := make([]types.Symptom, 0)
	if time_wait_sympton := GetTimeWaitSymptom(tree); time_wait_sympton != nil {
		symptoms = append(symptoms, time_wait_sympton)
	}

	if close_wait_sympton := GetCloseWaitSymptom(tree); close_wait_sympton != nil {
		symptoms = append(symptoms, close_wait_sympton)
	}

	if GetPortRangeSymptom := GetPortRangeSymptom(tree); GetPortRangeSymptom != nil {
		symptoms = append(symptoms, GetPortRangeSymptom)
	}
	return symptoms
}

func GetTimeWaitSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	netstat_status := metrics_tree.FindSection(netstat.NETSTAT_STATUS)
	time_wait_num, _ := strconv.Atoi(netstat_status.Find("TIME_WAIT"))
	if time_wait_num > MAX_TIME_OUT_CONNECTION {
		time_wait_symptom := &types.DefaultSymptom{
			Name:        TIME_WAIT_TOO_MUCH_SYMPTOM,
			Description: fmt.Sprintf("tcp connection state `TIME_WAIT` is too much, current amount is %d", time_wait_num),
			Advises: []types.Advise{
				&types.DefaultAdvise{
					Description: "`TIME_WAIT` means the client initiative close the connection and wait the stack to " +
						"recycle or reuse the connection, Maybe you use short connection in http client",
				},
				&types.DefaultAdvise{
					Description: "You can reuse tcp connection by set `keepalive` in http client,set `fastcgi_keep_conn` in php-fpm settings",
				},
				&types.DefaultAdvise{
					Description: "You can accelerate the `TIME_WAIT` connection recycle by sysctl: " +
						"sysclt -w net.ipv4.tcp_syncookies = 1;" +
						"sysclt -w net.ipv4.tcp_tw_reuse = 1;" +
						"sysclt -w net.ipv4.tcp_tw_recycle = 1;" +
						"sysclt -w net.ipv4.tcp_fin_timeout = 30",
				},
			},
		}
		return time_wait_symptom
	}
	return nil
}

func GetCloseWaitSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	netstat_status := metrics_tree.FindSection(netstat.NETSTAT_STATUS)
	close_wait_num, _ := strconv.Atoi(netstat_status.Find("CLOSE_WAIT"))
	if close_wait_num > MAX_CLOSE_WAIT_CONNECTION {
		time_wait_symptom := &types.DefaultSymptom{
			Name:        CLOSE_WAIT_TOO_MUCH_SYMPTOM,
			Description: fmt.Sprintf("tcp connection state `CLOSE_WAIT` is too much, current amount is %d", close_wait_num),
			Advises: []types.Advise{
				&types.DefaultAdvise{
					Description: "`CLOSE_WAIT` means some other application close the connection but you don't receive a fin pocket," +
						"You can check the api provider and close the connection timely",
				},
				&types.DefaultAdvise{
					Description: "`CLOSE_WAIT` could also occur when you client doesn't close response in http client.",
				},
			},
		}
		return time_wait_symptom
	}
	return nil
}

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
			ports_usage_symptom := &types.DefaultSymptom{
				Name:        PORTS_USAGE_TOO_MUCH_SYMPTOM,
				Description: fmt.Sprintf("Current system ports range is between %s, but ports total usage is %v", ports_total_range, ports_usage),
				Advises: []types.Advise{
					&types.DefaultAdvise{
						Description: "Ports Usage too much means connections is too much,Please check connection status is in normal status  " +
							"by `netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}' ` ",
					},
					&types.DefaultAdvise{
						Description: "Please check weather your application has too much 504 or 502 timeout in logs, " +
							"You can also increase port range by `sudo sysctl -w net.ipv4.ip_local_port_range=\"min_port_num max_port_num\"`",
					},
				},
			}
			return ports_usage_symptom
		}
	}
	return nil
}

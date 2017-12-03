package tcp

import (
	"fmt"
	"github.com/ringtail/snout/advisors"
	"github.com/ringtail/snout/collectors/netstat"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"strconv"
)

func init() {
	ta := &TcpAdvisor{}
	advisors.Add(ta)
}

var (
	TCP_ADVISOR                = "TCP_ADVISOR"
	TIME_WAIT_TOO_MUCH_SYMPTOM = "TIME_WAIT_TOO_MUCH"
)

const (
	MAX_TIME_OUT_CONNECTION = 200
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
	//kernel_settings := storage.InternalMetricsTree.FindSection(collectors.KERNEL_SETTINGS)
	symptoms := make([]types.Symptom, 0)
	netstat_status := storage.InternalMetricsTree.FindSection(netstat.NETSTAT_STATUS)
	time_wait_num, _ := strconv.Atoi(netstat_status.Find("TIME_WAIT"))
	if time_wait_num > MAX_TIME_OUT_CONNECTION {
		time_wait_symptom := &types.DefaultSymptom{
			Name:        TIME_WAIT_TOO_MUCH_SYMPTOM,
			Description: fmt.Sprintf("tcp connection state `TIME_WAIT` is too much, current amount is %d", time_wait_num),
			Advises: []types.Advise{
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
		symptoms = append(symptoms, time_wait_symptom)
	}
	return symptoms
}

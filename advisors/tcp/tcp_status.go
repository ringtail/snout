package tcp

import (
	"fmt"
	"github.com/ringtail/snout/collectors/netstat"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"strconv"
)

func GetTimeWaitSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	netstat_status := metrics_tree.FindSection(netstat.NETSTAT_STATUS)
	time_wait_num, _ := strconv.Atoi(netstat_status.Find("TIME_WAIT"))
	if time_wait_num > MAX_TIME_OUT_CONNECTION {
		desc := fmt.Sprintf("tcp connection state `TIME_WAIT` is too much, current amount is %d", time_wait_num)
		adviseDescs := []string{
			"`TIME_WAIT` means the client initiative close the connection and wait the stack to recycle or reuse the " +
				"connection, Maybe you use short connection in http client",

			"You can reuse tcp connection by set `keepalive` in http client,set `fastcgi_keep_conn` in php-fpm settings",

			"You can accelerate the `TIME_WAIT` connection recycle by sysctl: sysclt -w net.ipv4.tcp_syncookies = 1;" +
				"sysclt -w net.ipv4.tcp_tw_reuse = 1; sysclt -w net.ipv4.tcp_tw_recycle = 1; sysclt -w net.ipv4.tcp_fin_timeout = 30",
		}
		time_wait_symptom := types.CreateTextDefaultSymptom(TIME_WAIT_TOO_MUCH_SYMPTOM, desc, adviseDescs)
		return time_wait_symptom
	}
	return nil
}

func GetCloseWaitSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	netstat_status := metrics_tree.FindSection(netstat.NETSTAT_STATUS)
	close_wait_num, _ := strconv.Atoi(netstat_status.Find("CLOSE_WAIT"))
	if close_wait_num > MAX_CLOSE_WAIT_CONNECTION {
		desc := fmt.Sprintf("tcp connection state `CLOSE_WAIT` is too much, current amount is %d", close_wait_num)
		adviseDescs := []string{
			"`CLOSE_WAIT` means some other application close the connection but you don't receive a fin pocket," +
				"You can check the api provider and close the connection timely",

			"`CLOSE_WAIT` could also occur when you client doesn't close response in http client.",
		}
		close_wait_symptom := types.CreateTextDefaultSymptom(CLOSE_WAIT_TOO_MUCH_SYMPTOM, desc, adviseDescs)
		return close_wait_symptom
	}
	return nil
}

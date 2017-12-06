package netstat

import (
	"github.com/drael/GOnetstat"
	"github.com/ringtail/snout/collectors"
	"github.com/ringtail/snout/types"
	"strconv"
)

const (
	NETSTAT_STATUS = "NETSTAT_STATUS"
	PORTS_USAGE    = "PORTS_USAGE"
)

func init() {
	ns := &NetstatCollector{}
	collectors.Add(ns)
}

type NetstatCollector struct{}

func (nsc *NetstatCollector) Name() string {
	return NETSTAT_STATUS
}

func (nsc *NetstatCollector) Description() string {
	return ""
}

func (nsc *NetstatCollector) Gather() (types.MetricsSection, error) {
	tcp := GOnetstat.Tcp()
	tcp_connection_status := make(map[string]string)
	tcp_ports := make(map[int64]int)
	for _, p := range tcp {
		if tcp_connection_status[p.State] == "" {
			tcp_connection_status[p.State] = "1"
		} else {
			times, _ := strconv.Atoi(tcp_connection_status[p.State])
			tcp_connection_status[p.State] = strconv.Itoa(times + 1)
		}
		tcp_ports[p.Port] = tcp_ports[p.Port] + 1
	}

	tcp_connection_status [PORTS_USAGE] = strconv.Itoa(len(tcp_ports))

	return &types.DefaultMetricsSection{
		Name:    NETSTAT_STATUS,
		Metrics: tcp_connection_status,
	}, nil
}

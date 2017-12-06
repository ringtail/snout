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
	udp := GOnetstat.Udp()
	connection_status := make(map[string]string)
	ports := make(map[int64]int)
	for _, p := range tcp {
		if connection_status[p.State] == "" {
			connection_status[p.State] = "1"
		} else {
			times, _ := strconv.Atoi(connection_status[p.State])
			connection_status[p.State] = strconv.Itoa(times + 1)
		}
		ports[p.Port] = ports[p.Port] + 1
	}

	for _, u := range udp {
		ports[u.Port] = ports[u.Port] + 1
	}

	connection_status [PORTS_USAGE] = strconv.Itoa(len(ports))

	return &types.DefaultMetricsSection{
		Name:    NETSTAT_STATUS,
		Metrics: connection_status,
	}, nil
}

package tcp

import (
	"github.com/ringtail/snout/advisors"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
)

func init() {
	ta := &TcpAdvisor{}
	advisors.Add(ta)
}

var (
	TCP_ADVISOR                  = "TCP_ADVISOR"
	TIME_WAIT_TOO_MUCH_SYMPTOM   = "TIME_WAIT_TOO_MUCH"
	CLOSE_WAIT_TOO_MUCH_SYMPTOM  = "CLOSE_WAIT_TOO_MUCH"
	SYN_SENT_TOO_MUCH_SYMPTOM    = "SYN_SENT_TOO_MUCH"
	PORTS_USAGE_TOO_MUCH_SYMPTOM = "PORTS_USAGE_TOO_MUCH"
)

const (
	MAX_TIME_OUT_CONNECTION   = 20
	MAX_CLOSE_WAIT_CONNECTION = 20
	MAX_SYN_SENT_CONNECTION   = 10
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
	tree := storage.InternalMetricsTree
	symptoms := make([]types.Symptom, 0)
	if time_wait_sympton := GetTimeWaitSymptom(tree); time_wait_sympton != nil {
		symptoms = append(symptoms, time_wait_sympton)
	}

	if close_wait_sympton := GetCloseWaitSymptom(tree); close_wait_sympton != nil {
		symptoms = append(symptoms, close_wait_sympton)
	}

	if port_range_symptom := GetPortRangeSymptom(tree); port_range_symptom != nil {
		symptoms = append(symptoms, port_range_symptom)
	}

	if syn_sent_symptom := GetSynSentSymptom(tree); syn_sent_symptom != nil {
		symptoms = append(symptoms, syn_sent_symptom)
	}
	return symptoms
}

package core

import (
	log "github.com/Sirupsen/logrus"

	"github.com/ringtail/snout/advisors"
	"github.com/ringtail/snout/collectors"
	"github.com/ringtail/snout/types"
)

type Plugins interface {
	Empty() bool
	Start()
}

type Snout struct {
	CollectorManager types.Manager
	AdvisorsManager  types.Manager
}

func (st *Snout) Load() {
	st.CollectorManager = collectors.Cm
	st.AdvisorsManager = advisors.Am
}

func (st *Snout) Run() {
	st.Load()
	if st.CollectorManager.Empty() {
		log.Errorf("Failed to load any collectors, advisors or resolvers.")
		return
	}
	st.CollectorManager.Start()
	st.AdvisorsManager.Start()
}

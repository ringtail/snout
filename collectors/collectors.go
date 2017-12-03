package collectors

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/types"
	"sync"
)

type CollectorInterface interface {
	Name() string
	Description() string
	Gather() (types.MetricsSection, error)
}

type CollectionManager struct {
	Collectors map[string]CollectorInterface
}

func (cm *CollectionManager) Empty() bool {
	return false
}

func (cm *CollectionManager) Add(ci CollectorInterface) {
	if ci.Name() == "" {
		return
	}
	cm.Collectors[ci.Name()] = ci
}

func (cm *CollectionManager) Find(name string) CollectorInterface {
	return cm.Collectors[name]
}

func (cm *CollectionManager) Start(names ...string) {
	var wg sync.WaitGroup
	if len(names) == 0 {
		for index, _ := range cm.Collectors {
			ci := cm.Collectors[index]
			wg.Add(1)
			go func(collectorInterface CollectorInterface) {
				defer wg.Done()
				ms, err := ci.Gather()
				if err != nil {
					log.Warnf("[%s] Failed to gather metrics,because of %v", ci.Name(), err.Error())
					return
				}
				storage.InternalMetricsTree.AddSection(ms)
			}(ci)
		}
	} else {
		for _, name := range names {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ci := cm.Find(name)
				if ci != nil {
					ms, err := ci.Gather()
					if err != nil {
						log.Warnf("[%s] Failed to gather metrics,because of %v", ci.Name(), err.Error())
						return
					}
					storage.InternalMetricsTree.AddSection(ms)
				}
			}()
		}
	}
	wg.Wait()
}

var Cm *CollectionManager

func init() {
	Cm = &CollectionManager{}
	Cm.Collectors = make(map[string]CollectorInterface, 0)
}

func Add(ci CollectorInterface) {
	Cm.Add(ci)
}

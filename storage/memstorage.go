package storage

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ringtail/snout/types"
)

type MetricsTree struct {
	MetricsSection map[string]types.MetricsSection
}

func (mt *MetricsTree) AddSection(section types.MetricsSection) {
	name := section.GetName()
	mt.MetricsSection[name] = section
}

func (mt *MetricsTree) FindSection(name string) types.MetricsSection {
	return mt.MetricsSection[name]
}

func (mt *MetricsTree) DumpAll() {
	for key, value := range mt.MetricsSection {
		log.Debugf("================== Metric Section %s Begin ==================", key)
		for map_key, map_value := range value.List() {
			log.Debugf("%s:%s", map_key, map_value)
		}
		log.Debugf("================== Metric Section %s End ==================", key)
	}
}

type SymptomStorage struct {
}

var (
	InternalMetricsTree *MetricsTree
	InternalSymptom     *SymptomStorage
)

func init() {
	InternalMetricsTree = &MetricsTree{}
	InternalMetricsTree.MetricsSection = make(map[string]types.MetricsSection)
}

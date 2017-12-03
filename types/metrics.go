package types

type MetricsSection interface {
	GetName() string
	List() map[string]string
	Find(name string) string
}

type DefaultMetricsSection struct {
	Name    string
	Metrics map[string]string
}

func (dms *DefaultMetricsSection) GetName() string {
	return dms.Name
}
func (dms *DefaultMetricsSection) List() map[string]string {
	return dms.Metrics
}

func (dms *DefaultMetricsSection) Find(name string) string {
	return dms.Metrics[name]
}

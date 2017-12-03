package advisors

import (
	"github.com/ringtail/snout/types"
)

type AdvisorInterface interface {
	Name() string
	Description() string
	Advise() []types.Symptom
}

type AdvisorsManager struct {
	Advisors map[string]AdvisorInterface
}

func (am *AdvisorsManager) Find(name string) AdvisorInterface {
	return am.Advisors[name]
}

func (am *AdvisorsManager) Add(name string, adi AdvisorInterface) {
	am.Advisors[name] = adi
}

func (am *AdvisorsManager) Empty() bool {
	return false
}

func (am *AdvisorsManager) Start(names ...string) {
	dra := types.DefaultDiagnosticReport{}
	if len(names) == 0 {
		for _, advisor := range am.Advisors {
			symptoms := advisor.Advise()
			dra.Add(symptoms)
		}
	} else {
		for _, name := range names {
			advisor := am.Find(name)
			symptoms := advisor.Advise()
			dra.Add(symptoms)
		}
	}
	dra.Print()
}

var Am *AdvisorsManager

func init() {
	Am = &AdvisorsManager{}
	Am.Advisors = make(map[string]AdvisorInterface)
}

func Add(ad AdvisorInterface) {
	name := ad.Name()
	if name == "" {
		return
	}
	Am.Add(name, ad)
}

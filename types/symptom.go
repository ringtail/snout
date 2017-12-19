package types

import (
	"github.com/olekukonko/tablewriter"
	"github.com/ringtail/snout/resolvers"
	"os"
	"fmt"
)

/**
A DiagnosticReport may contains several symptom
A symptom
*/

type DiagnosticReport interface {
	GetSymptom() []Symptom
	Print()
}

type Symptom interface {
	GetDescription() string
	GetName() string
	GetAdvises() []Advise
	GetAdviseDescriptions() []string
}

type Advise interface {
	GetDescription() string
	GetResolvers() []resolvers.Resolver
}

type DefaultAdvise struct {
	Description string
	Resolvers   []resolvers.Resolver
}

func (da *DefaultAdvise) GetDescription() string {
	return da.Description
}

func (da *DefaultAdvise) GetResolvers() []resolvers.Resolver {
	return da.Resolvers
}

type DefaultDiagnosticReport struct {
	Symptom []Symptom
}

func (ddr *DefaultDiagnosticReport) Add(syms []Symptom) {
	for _, sym := range syms {
		ddr.Symptom = append(ddr.Symptom, sym)
	}
}

func (ddr *DefaultDiagnosticReport) GetSymptom() []Symptom {
	return ddr.Symptom
}

func (ddr *DefaultDiagnosticReport) Print() {
	if len(ddr.GetSymptom()) == 0 {
		fmt.Println("Snout doesn't sniff any bad smell about performance in your system")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	//table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Symptom", "Description", "Advises"})

	for _, symptom := range ddr.GetSymptom() {
		data := make([][]string, 0)
		for _, adviseDescription := range symptom.GetAdviseDescriptions() {
			data = append(data, []string{symptom.GetName(), symptom.GetDescription(), adviseDescription})
		}

		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.SetColMinWidth(2, 80)
		table.AppendBulk(data) // Add Bulk Data

	}
	table.Render()
}

type DefaultSymptom struct {
	Name        string
	Description string
	Advises     []Advise
}

func (ds *DefaultSymptom) GetName() string {
	return ds.Name
}
func (ds *DefaultSymptom) GetDescription() string {
	return ds.Description
}

func (ds *DefaultSymptom) GetAdvises() []Advise {
	return ds.Advises
}

func (ds *DefaultSymptom) GetAdviseDescriptions() []string {
	descriptions := make([]string, 0)
	for _, advise := range ds.Advises {
		descriptions = append(descriptions, advise.GetDescription())
	}
	return descriptions
}

func CreateTextDefaultSymptom(name, desc string, adviseDescriptions []string) Symptom {
	symptom := &DefaultSymptom{
		Name:        name,
		Description: desc,
		Advises:     make([]Advise, 0),
	}
	for _, adviseDesc := range adviseDescriptions {
		symptom.Advises = append(symptom.Advises, &DefaultAdvise{
			Description: adviseDesc,
		})
	}
	return symptom
}

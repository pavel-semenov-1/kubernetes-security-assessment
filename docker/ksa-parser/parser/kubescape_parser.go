package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type KubescapeParser struct {
	data  KubescapeData
	mutex sync.Mutex
}

type KubescapeData struct {
	Details SummaryDetails `json:"summaryDetails"`
}

type SummaryDetails struct {
	Controls map[string]KubescapeControl `json:"controls"`
}

type KubescapeControl struct {
	Id       string            `json:"controlID"`
	Name     string            `json:"name"`
	Status   string            `json:"status"`
	Category KubescapeCategory `json:"category"`
}

type KubescapeCategory struct {
	Name string `json:"name"`
}

func NewKubescapeParser() *KubescapeParser {
	return &KubescapeParser{}
}

func (p *KubescapeParser) Parse(filePath string) ([]Vulnerability, []Misconfiguration, error) {
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	var result KubescapeData
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	p.mutex.Lock()
	p.data = result
	p.mutex.Unlock()

	return p.GetVulnerabilities(), p.GetMisconfigurations(), nil
}

func (p *KubescapeParser) GetResults() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.data
}

func (p *KubescapeParser) GetVulnerabilities() []Vulnerability {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return []Vulnerability{}
}

func (p *KubescapeParser) GetMisconfigurations() []Misconfiguration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var misconfigurations []Misconfiguration
	for _, ctrl := range p.data.Details.Controls {
		if ctrl.Status == "skipped" {
			continue
		}
		var status string
		if ctrl.Status == "passed" {
			status = "PASS"
		} else if ctrl.Status == "failed" {
			status = "FAIL"
		}
		misconfigurations = append(misconfigurations, Misconfiguration{
			Type:               ctrl.Category.Name,
			MisconfigurationID: ctrl.Id,
			Title:              ctrl.Name,
			Description:        ctrl.Name,
			Resolution:         "No remediation provided.",
			Severity:           "HIGH",
			Target:             "Security Control",
			Status:             status,
		})
	}
	return misconfigurations
}

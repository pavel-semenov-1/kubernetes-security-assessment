package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type KubeBenchParser struct {
	data  KubeBenchResult
	mutex sync.Mutex
}

type KubeBenchResult struct {
	Controls []Control `json:"Controls"`
	Totals   Totals    `json:"Totals"`
}

type Control struct {
	Id    string `json:"id"`
	Text  string `json:"text"`
	Tests []Test `json:"tests"`
}

type Totals struct {
	TotalPass int `json:"total_pass"`
	TotalFail int `json:"total_fail"`
	TotalWarn int `json:"total_warn"`
	TotalInfo int `json:"total_info"`
}

type Test struct {
	Pass              int                `json:"pass"`
	Fail              int                `json:"fail"`
	Warn              int                `json:"warn"`
	Info              int                `json:"info"`
	Desc              string             `json:"desc"`
	Results           []interface{}      `json:"results"`                     // original from the report
	Misconfigurations []Misconfiguration `json:"misconfigurations,omitempty"` // normalized
}

func NewKubeBenchParser() *KubeBenchParser {
	return &KubeBenchParser{}
}

func (p *KubeBenchParser) Parse(filePath string) ([]Vulnerability, []Misconfiguration, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	var result KubeBenchResult
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Normalize Misconfigurations
	for ctrlIndex, control := range result.Controls {
		for testIndex, t := range control.Tests {
			var normalizedResults []Misconfiguration
			for _, raw := range t.Results {
				if resultMap, ok := raw.(map[string]interface{}); ok {
					normalizedResults = append(normalizedResults, Misconfiguration{
						Type:        resultMap["type"].(string),
						ID:          resultMap["test_number"].(string),
						Title:       t.Desc,
						Description: resultMap["test_desc"].(string),
						Resolution:  resultMap["remediation"].(string),
						Severity:    resultMap["status"].(string),
					})
				}
			}
			// Replace Results with normalized misconfigurations
			result.Controls[ctrlIndex].Tests[testIndex].Misconfigurations = normalizedResults
		}
	}

	p.data = result
	return p.GetVulnerabilities(), p.GetMisconfigurations(), nil
}

func (p *KubeBenchParser) GetResults() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.data
}

func (p *KubeBenchParser) GetVulnerabilities() []Vulnerability {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return []Vulnerability{}
}

func (p *KubeBenchParser) GetMisconfigurations() []Misconfiguration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var misconfigurations []Misconfiguration
	for _, ctrl := range p.data.Controls {
		for _, test := range ctrl.Tests {
			if test.Misconfigurations != nil {
				misconfigurations = append(misconfigurations, test.Misconfigurations...)
			}
		}
	}
	return misconfigurations
}

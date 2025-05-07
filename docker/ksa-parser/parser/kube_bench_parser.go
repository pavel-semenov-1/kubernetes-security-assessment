package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type KubeBenchParser struct {
	data  []KubeBenchControl
	mutex sync.Mutex
}

type KubeBenchControl struct {
	Id    string `json:"id"`
	Text  string `json:"text"`
	Tests []Test `json:"tests"`
}

type Test struct {
	Desc    string            `json:"desc"`
	Results []KubeBenchResult `json:"results"`
}

type KubeBenchResult struct {
	TestNumber  string `json:"test_number"`
	Status      string `json:"status"`
	TestDesc    string `json:"test_desc"`
	Remediation string `json:"remediation"`
}

func NewKubeBenchParser() *KubeBenchParser {
	return &KubeBenchParser{}
}

func (p *KubeBenchParser) Parse(filePath string) ([]Vulnerability, []Misconfiguration, error) {
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	var result []KubeBenchControl
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	p.mutex.Lock()
	p.data = result
	p.mutex.Unlock()

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
	for _, ctrl := range p.data {
		for _, test := range ctrl.Tests {
			for _, result := range test.Results {
				status := result.Status
				if status == "WARN" {
					status = "MANUAL"
				}
				misconfigurations = append(misconfigurations, Misconfiguration{
					Type:               "Kubernetes security check",
					MisconfigurationID: result.TestNumber,
					Title:              fmt.Sprintf("%s. %s - %s", result.TestNumber, ctrl.Text, test.Desc),
					Description:        result.TestDesc,
					Resolution:         result.Remediation,
					Severity:           "HIGH",
					Target:             "Security Check",
					Status:             status,
				})
			}
		}
	}
	return misconfigurations
}

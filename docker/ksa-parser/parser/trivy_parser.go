package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type TrivyParser struct {
	data  TrivyResult
	mutex sync.Mutex
}

type TrivyResult struct {
	ClusterName string     `json:"ClusterName"`
	Resources   []Resource `json:"Resources"`
}

type Resource struct {
	Namespace string   `json:"Namespace"`
	Kind      string   `json:"Kind"`
	Name      string   `json:"Name"`
	Results   []Result `json:"Results"`
}

type Result struct {
	Target            string             `json:"Target"`
	Class             string             `json:"Class"`
	Type              string             `json:"Type"`
	Vulnerabilities   []Vulnerability    `json:"Vulnerabilities,omitempty"`
	MisconfSummary    MisconfSummary     `json:"MisconfSummary,omitempty"`
	Misconfigurations []Misconfiguration `json:"Misconfigurations,omitempty"`
}

type MisconfSummary struct {
	Successes  int `json:"Successes"`
	Failures   int `json:"Failures"`
	Exceptions int `json:"Exceptions"`
}

func NewTrivyParser() *TrivyParser {
	return &TrivyParser{}
}

func (p *TrivyParser) Parse(filePath string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var result TrivyResult
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	p.data = result
	return nil
}

func (p *TrivyParser) GetResults() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.data
}

func (p *TrivyParser) GetVulnerabilities(namespace *string, severity *string) []Vulnerability {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var vulnerabilities []Vulnerability
	for _, res := range p.data.Resources {
		if namespace == nil || *namespace == "" || *namespace == res.Namespace {
			for _, rslt := range res.Results {
				if rslt.Vulnerabilities != nil {
					if severity != nil && *severity != "" {
						severityFilter := func(v Vulnerability) bool { return v.Severity == *severity }
						vulnerabilities = append(vulnerabilities, filterBy(rslt.Vulnerabilities, severityFilter)...)
					} else {
						vulnerabilities = append(vulnerabilities, rslt.Vulnerabilities...)
					}
				}
			}
		}
	}
	return vulnerabilities
}

func (p *TrivyParser) GetMisconfigurations(namespace *string, severity *string) []Misconfiguration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var misconfigurations []Misconfiguration
	for _, res := range p.data.Resources {
		if namespace == nil || *namespace == "" || *namespace == res.Namespace {
			for _, rslt := range res.Results {
				if rslt.Misconfigurations != nil {
					if severity != nil && *severity != "" {
						severityFilter := func(v Misconfiguration) bool { return v.Severity == *severity }
						misconfigurations = append(misconfigurations, filterBy(rslt.Misconfigurations, severityFilter)...)
					} else {
						misconfigurations = append(misconfigurations, rslt.Misconfigurations...)
					}
				}
			}
		}
	}
	return misconfigurations
}

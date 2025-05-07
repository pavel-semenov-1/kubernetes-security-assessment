package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type TrivyParser struct {
	data  TrivyReport
	mutex sync.Mutex
}

type TrivyReport struct {
	ClusterName string          `json:"ClusterName"`
	Resources   []TrivyResource `json:"Resources"`
}

type TrivyResource struct {
	Namespace string        `json:"Namespace"`
	Kind      string        `json:"Kind"`
	Name      string        `json:"Name"`
	Results   []TrivyResult `json:"Results"`
}

type TrivyResult struct {
	Target            string                  `json:"Target"`
	Class             string                  `json:"Class"`
	Type              string                  `json:"Type"`
	Vulnerabilities   []Vulnerability         `json:"Vulnerabilities,omitempty"`
	MisconfSummary    MisconfSummary          `json:"MisconfSummary,omitempty"`
	Misconfigurations []TrivyMisconfiguration `json:"Misconfigurations,omitempty"`
}

type TrivyMisconfiguration struct {
	Type               string `json:"Type"`
	MisconfigurationID string `json:"ID"`
	Title              string `json:"Title"`
	Description        string `json:"Description"`
	Resolution         string `json:"Resolution"`
	Severity           string `json:"Severity"`
	Target             string `json:"Target"`
	Status             string `json:"Status"`
}

type MisconfSummary struct {
	Successes  int `json:"Successes"`
	Failures   int `json:"Failures"`
	Exceptions int `json:"Exceptions"`
}

func NewTrivyParser() *TrivyParser {
	return &TrivyParser{}
}

func (p *TrivyParser) Parse(filePath string) ([]Vulnerability, []Misconfiguration, error) {
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	var result TrivyReport
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	for resourceIndex, res := range result.Resources {
		for resultIndex, rslt := range res.Results {
			idExtractorVuln := func(v Vulnerability) string { return v.VulnerabilityID }
			result.Resources[resourceIndex].Results[resultIndex].Vulnerabilities = RemoveDuplicates(rslt.Vulnerabilities, idExtractorVuln)
			idExtractorMisc := func(m TrivyMisconfiguration) string { return m.MisconfigurationID }
			result.Resources[resourceIndex].Results[resultIndex].Misconfigurations = RemoveDuplicates(rslt.Misconfigurations, idExtractorMisc)
			for vulnIndex, _ := range result.Resources[resourceIndex].Results[resultIndex].Vulnerabilities {
				result.Resources[resourceIndex].Results[resultIndex].Vulnerabilities[vulnIndex].Target = rslt.Target
			}
			for misconfIndex, _ := range result.Resources[resourceIndex].Results[resultIndex].Misconfigurations {
				result.Resources[resourceIndex].Results[resultIndex].Misconfigurations[misconfIndex].Target = rslt.Target
			}
		}
	}

	p.mutex.Lock()
	p.data = result
	p.mutex.Unlock()

	return p.GetVulnerabilities(), p.GetMisconfigurations(), nil
}

func (p *TrivyParser) GetResults() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.data
}

func (p *TrivyParser) GetVulnerabilities() []Vulnerability {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var vulnerabilities []Vulnerability
	for _, res := range p.data.Resources {
		for _, rslt := range res.Results {
			if rslt.Vulnerabilities != nil {
				vulnerabilities = append(vulnerabilities, rslt.Vulnerabilities...)
			}
		}
	}
	return vulnerabilities
}

func (p *TrivyParser) GetMisconfigurations() []Misconfiguration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var misconfigurations []Misconfiguration
	for _, res := range p.data.Resources {
		for _, rslt := range res.Results {
			for _, misconf := range rslt.Misconfigurations {
				misconfigurations = append(misconfigurations, Misconfiguration{
					Type:               misconf.Type,
					MisconfigurationID: misconf.MisconfigurationID,
					Title:              misconf.Title,
					Description:        misconf.Description,
					Resolution:         misconf.Resolution,
					Severity:           misconf.Severity,
					Target:             misconf.Target,
					Status:             misconf.Status,
				})
			}
		}
	}
	return misconfigurations
}

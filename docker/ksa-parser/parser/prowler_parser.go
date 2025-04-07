package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type ProwlerParser struct {
	data  []ProwlerReport
	mutex sync.Mutex
}

type ProwlerReport struct {
	FindingInfo FindingInfo       `json:"finding_info"`
	Severity    string            `json:"severity"`
	Remediation Remediation       `json:"remediation"`
	Resources   []ProwlerResource `json:"resources"`
	TypeName    string            `json:"type_name"`
}

type FindingInfo struct {
	Desc  string `json:"desc"`
	Title string `json:"title"`
	UID   string `json:"uid"`
}

type Remediation struct {
	Desc string `json:"desc"`
}

type ProwlerResource struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func NewProwlerParser() *ProwlerParser {
	return &ProwlerParser{}
}

func (p *ProwlerParser) Parse(filePath string) ([]Vulnerability, []Misconfiguration, error) {
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	var result []ProwlerReport
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	p.mutex.Lock()
	p.data = result
	p.mutex.Unlock()

	return p.GetVulnerabilities(), p.GetMisconfigurations(), nil
}

func (p *ProwlerParser) GetResults() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.data
}

func (p *ProwlerParser) GetVulnerabilities() []Vulnerability {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var vulnerabilities []Vulnerability
	return vulnerabilities
}

func (p *ProwlerParser) GetMisconfigurations() []Misconfiguration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var misconfigurations []Misconfiguration
	for _, res := range p.data {
		var target string
		for rsIdx, rs := range res.Resources {
			if rsIdx != 0 {
				target = fmt.Sprintf("%s, ", target)
			}
			target = fmt.Sprintf("%s%s/%s", target, rs.Namespace, rs.Name)
		}
		misconfigurations = append(misconfigurations, Misconfiguration{
			Type:        res.TypeName,
			ID:          res.FindingInfo.UID,
			Title:       res.FindingInfo.Title,
			Description: res.FindingInfo.Desc,
			Resolution:  res.Remediation.Desc,
			Severity:    res.Severity,
			Target:      target,
		})
	}
	return misconfigurations
}

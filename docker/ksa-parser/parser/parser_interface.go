package parser

// Parser defines a common interface for all security scanner parsers
type Parser interface {
	Parse(filePath string) ([]Vulnerability, []Misconfiguration, error)
	GetResults() interface{}
	GetVulnerabilities() []Vulnerability
	GetMisconfigurations() []Misconfiguration
}

type Vulnerability struct {
	ID               int    `json:"ID"`
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
	Severity         string `json:"Severity"`
	Target           string `json:"Target"`
}

type Misconfiguration struct {
	ID                 int    `json:"ID"`
	Type               string `json:"Type"`
	MisconfigurationID string `json:"MisconfigurationID"`
	Title              string `json:"Title"`
	Description        string `json:"Description"`
	Resolution         string `json:"Resolution"`
	Severity           string `json:"Severity"`
	Target             string `json:"Target"`
	Status             string `json:"Status"`
}

func RemoveDuplicates[T any](ss []T, extractor func(T) string) (ret []T) {
	ids := make(map[string]bool)
	for _, s := range ss {
		if !ids[extractor(s)] {
			ids[extractor(s)] = true
			ret = append(ret, s)
		}
	}
	return
}

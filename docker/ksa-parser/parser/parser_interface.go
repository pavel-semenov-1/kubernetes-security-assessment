package parser

// Parser defines a common interface for all security scanner parsers
type Parser interface {
	Parse(filePath string) error
	GetResults() interface{}
	GetVulnerabilities(namespace *string, severity *string) []Vulnerability
	GetMisconfigurations(namespace *string, severity *string) []Misconfiguration
}

type Vulnerability struct {
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
	Type        string `json:"Type"`
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Resolution  string `json:"Resolution"`
	Severity    string `json:"Severity"`
	Target      string `json:"Target"`
}

func filterBy[T any](ss []T, filter func(T) bool) (ret []T) {
	for _, s := range ss {
		if filter(s) {
			ret = append(ret, s)
		}
	}
	return
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

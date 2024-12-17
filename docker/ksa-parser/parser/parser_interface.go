package parser

// Parser defines a common interface for all security scanner parsers
type Parser interface {
	Parse(filePath string) error // Parse the input file
	GetResults() interface{}     // Retrieve the parsed data
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
}

type Misconfiguration struct {
	Type        string `json:"Type"`
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Resolution  string `json:"Resolution"`
	Severity    string `json:"Severity"`
}

func filterBy[T any](ss []T, filter func(T) bool) (ret []T) {
	for _, s := range ss {
		if filter(s) {
			ret = append(ret, s)
		}
	}
	return
}

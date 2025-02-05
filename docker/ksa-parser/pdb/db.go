package pdb

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"ksa-parser/parser"
)

type DB struct {
	connection *sql.DB
}

type Report struct {
	Id       int    `json:"ID"`
	Filename string `json:"Filename"`
}

func NewDB(connStr string) (*DB, error) {
	con, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("unable to open DB: %v", err))
	}

	return &DB{connection: con}, nil
}

func (db *DB) Close() {
	err := db.connection.Close()
	if err != nil {
		panic(fmt.Errorf("unable to close DB: %v", err))
	}
}

func (db *DB) Parsed(reportId string) (bool, error) {
	rows, err := db.connection.Query("SELECT parsed FROM report where id=$1", reportId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		var parsed bool
		err := rows.Scan(&parsed)
		if err != nil {
			panic(err)
		}
		return parsed, nil
	}

	return false, errors.New("no report found")
}

func (db *DB) GetReportFilename(reportId string) string {
	rows, err := db.connection.Query("SELECT filename FROM report where id=$1", reportId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var filename string
		err := rows.Scan(&filename)
		if err != nil {
			panic(err)
		}
		return filename
	}

	panic("no report found")
}

func (db *DB) Populate(reportId string, vuln []parser.Vulnerability, misc []parser.Misconfiguration) error {
	vulnStmt, err := db.connection.Prepare("INSERT INTO vulnerability (report_id, vid, pkg_name, installed_version, fixed_version, title, description, severity, target) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		return err
	}
	defer vulnStmt.Close()

	for _, v := range vuln {
		_, err = vulnStmt.Exec(reportId, v.VulnerabilityID, v.PkgName, v.InstalledVersion, v.FixedVersion, v.Title, v.Description, v.Severity, v.Target)
		if err != nil {
			return err
		}
	}

	miscStmt, err := db.connection.Prepare("INSERT INTO misconfiguration (report_id, mid, type, title, description, resolution, severity, target) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		return err
	}
	defer miscStmt.Close()

	for _, m := range misc {
		_, err = miscStmt.Exec(reportId, m.ID, m.Type, m.Title, m.Description, m.Resolution, m.Severity, m.Target)
		if err != nil {
			return err
		}
	}

	_, err = db.connection.Exec("UPDATE report SET parsed=true WHERE id=$1", reportId)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetVulnerabilities(reportId string, searchTerm string) ([]parser.Vulnerability, error) {
	rows, err := db.connection.Query(`SELECT vid, pkg_name, installed_version, fixed_version, title, description, severity, target 
													FROM vulnerability WHERE report_id = $1 
											        		AND ($2 = '' OR search_vector @@ plainto_tsquery($2))`, reportId, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vulnerabilities := make([]parser.Vulnerability, 0)
	for rows.Next() {
		var v parser.Vulnerability
		if err := rows.Scan(&v.VulnerabilityID, &v.PkgName, &v.InstalledVersion, &v.FixedVersion,
			&v.Title, &v.Description, &v.Severity, &v.Target); err != nil {
			return nil, err
		}

		vulnerabilities = append(vulnerabilities, v)
	}
	return vulnerabilities, nil
}

func (db *DB) GetMisconfigurations(reportId string, searchTerm string) ([]parser.Misconfiguration, error) {
	rows, err := db.connection.Query(`SELECT type, mid, title, description, resolution, severity, target 
													FROM misconfiguration WHERE report_id = $1 
											        		AND ($2 = '' OR search_vector @@ plainto_tsquery($2))`, reportId, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	misconfigurations := make([]parser.Misconfiguration, 0)
	for rows.Next() {
		var m parser.Misconfiguration
		if err := rows.Scan(&m.Type, &m.ID, &m.Title, &m.Description, &m.Resolution, &m.Severity, &m.Target); err != nil {
			return nil, err
		}

		misconfigurations = append(misconfigurations, m)
	}
	return misconfigurations, nil
}

func (db *DB) GetReports(scanner string) ([]Report, error) {
	reports := []Report{}
	rows, err := db.connection.Query("SELECT id FROM scanner where name=$1", scanner)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		rows, err = db.connection.Query("SELECT id, filename FROM report where scanner_id=$1", id)
		if err != nil {
			return nil, fmt.Errorf("query failed: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var filename string
			err = rows.Scan(&id, &filename)
			if err != nil {
				return nil, fmt.Errorf("row scan failed: %v", err)
			}

			reports = append(reports, Report{Id: id, Filename: filename})
		}
	}

	return reports, nil
}

func (db *DB) GetLastParsedReportId(scanner string) (int, error) {
	rows, err := db.connection.Query("SELECT id FROM scanner where name=$1", scanner)
	if err != nil {
		return 0, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("row scan failed: %v", err)
		}
		rows, err = db.connection.Query("SELECT id FROM report where scanner_id=$1 and parsed=true ORDER BY generated_at DESC LIMIT 1", id)
		if err != nil {
			return 0, fmt.Errorf("query failed: %v", err)
		}
		defer rows.Close()

		if rows.Next() {
			var reportId int
			err = rows.Scan(&reportId)
			if err != nil {
				return 0, fmt.Errorf("row scan failed: %v", err)
			}

			return reportId, nil
		}
	}

	return 0, nil
}

func (db *DB) GetVulnerabilityStatistics(reportId int) (map[string]int, error) {
	results := map[string]int{
		"CRITICAL": 0,
		"HIGH":     0,
		"MEDIUM":   0,
		"LOW":      0,
	}

	rows, err := db.connection.Query("SELECT severity, count(*) FROM vulnerability where report_id=$1 GROUP BY severity;", reportId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int
		err = rows.Scan(&severity, &count)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		results[severity] = count
	}

	return results, nil
}

func (db *DB) GetMisconfigurationStatistics(reportId int) (map[string]int, error) {
	results := map[string]int{
		"CRITICAL": 0,
		"HIGH":     0,
		"MEDIUM":   0,
		"LOW":      0,
	}

	rows, err := db.connection.Query("SELECT severity, count(*) FROM misconfiguration where report_id=$1 GROUP BY severity;", reportId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int
		err = rows.Scan(&severity, &count)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		results[severity] = count
	}

	return results, nil
}

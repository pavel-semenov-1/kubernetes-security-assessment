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

func (db *DB) GetVulnerabilities(reportId string, severity string) ([]parser.Vulnerability, error) {
	vulnerabilities := []parser.Vulnerability{}
	rows, err := db.connection.Query("SELECT vid, pkg_name, installed_version, fixed_version, title, description, target FROM vulnerability where report_id=$1 and severity=$2", reportId, severity)
	if err != nil {
		return vulnerabilities, err
	}
	defer rows.Close()

	for rows.Next() {
		var vid string
		var pkgName string
		var installedVersion string
		var fixedVersion string
		var title string
		var description string
		var target string
		err := rows.Scan(&vid, &pkgName, &installedVersion, &fixedVersion, &title, &description, &target)
		if err != nil {
			return vulnerabilities, err
		}
		vulnerabilities = append(vulnerabilities, parser.Vulnerability{
			VulnerabilityID:  vid,
			PkgName:          pkgName,
			InstalledVersion: installedVersion,
			FixedVersion:     fixedVersion,
			Title:            title,
			Description:      description,
			Severity:         severity,
			Target:           target,
		})
	}
	return vulnerabilities, nil
}

func (db *DB) GetMisconfigurations(reportId string, severity string) ([]parser.Misconfiguration, error) {
	misconfigurations := []parser.Misconfiguration{}
	rows, err := db.connection.Query("SELECT mid, type, title, description, resolution, target FROM misconfiguration where report_id=$1 and severity=$2", reportId, severity)
	if err != nil {
		return misconfigurations, err
	}
	defer rows.Close()

	for rows.Next() {
		var mid string
		var ttype string
		var title string
		var description string
		var resolution string
		var target string
		err := rows.Scan(&mid, &ttype, &title, &description, &resolution, &target)
		if err != nil {
			return misconfigurations, err
		}
		misconfigurations = append(misconfigurations, parser.Misconfiguration{
			ID:          mid,
			Type:        ttype,
			Title:       title,
			Description: description,
			Resolution:  resolution,
			Severity:    severity,
			Target:      target,
		})
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

	stmt, err := db.connection.Prepare("SELECT count(*) FROM vulnerability where report_id=$1 and severity=$2")
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer stmt.Close()

	for svrt := range results {
		rows, err := stmt.Query(reportId, svrt)
		if err != nil {
			return nil, fmt.Errorf("query failed: %v", err)
		}
		defer rows.Close()

		if rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				return nil, fmt.Errorf("row scan failed: %v", err)
			}
			results[svrt] = count
		}
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

	stmt, err := db.connection.Prepare("SELECT count(*) FROM misconfiguration where report_id=$1 and severity=$2")
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer stmt.Close()

	for svrt := range results {
		rows, err := stmt.Query(reportId, svrt)
		if err != nil {
			return nil, fmt.Errorf("query failed: %v", err)
		}
		defer rows.Close()

		if rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				return nil, fmt.Errorf("row scan failed: %v", err)
			}
			results[svrt] = count
		}
	}

	return results, nil
}

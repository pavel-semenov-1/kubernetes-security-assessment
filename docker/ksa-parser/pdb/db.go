package pdb

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"ksa-parser/parser"
	"strings"
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

func (db *DB) GetScannerNames() ([]string, error) {
	rows, err := db.connection.Query("SELECT name FROM scanner ORDER BY name")
	if err != nil {
		return nil, err
	}

	scannerNames := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		scannerNames = append(scannerNames, name)
	}

	return scannerNames, nil
}

func (db *DB) Populate(reportId string, vuln []parser.Vulnerability, misc []parser.Misconfiguration) error {
	stmt, err := db.connection.Prepare(`
	  INSERT INTO vulnerability (
		  report_id, vid, pkg_name, installed_version,
		  fixed_version, title, description, severity,
		  target
		)
		SELECT 
		  $1, $2, $3, $4,
		  $5, $6, $7, $8,
		  $9
		FROM (
		  SELECT $2::text AS vid, $9::text AS target
		) AS new_vuln
		LEFT JOIN (
		  SELECT DISTINCT ON (vid, target) vid, target
		  FROM vulnerability
		  ORDER BY vid, target
		) v ON v.vid = new_vuln.vid AND v.target = new_vuln.target`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range vuln {
		_, err = stmt.Exec(
			reportId,
			v.VulnerabilityID,
			v.PkgName,
			v.InstalledVersion,
			v.FixedVersion,
			v.Title,
			v.Description,
			v.Severity,
			v.Target,
		)
		if err != nil {
			return err
		}
	}

	stmt, err = db.connection.Prepare(`INSERT INTO misconfiguration (
		  report_id, mid, type, title,
		  description, resolution, severity,
		  target, status
		)
		SELECT 
		  $1, $2, $3, $4,
		  $5, $6, $7,
		  $8,
		  CASE
			WHEN m.status = 'RESOLVED' THEN m.status
			ELSE COALESCE($9, 'FAIL')
		  END
		FROM (
		  SELECT $2::text AS mid, $8::text AS target
		) AS new_misc
		LEFT JOIN (
		  SELECT DISTINCT ON (mid, target) mid, target, status
		  FROM misconfiguration
		  ORDER BY mid, target
		) m ON m.mid = new_misc.mid AND m.target = new_misc.target
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range misc {
		_, err = stmt.Exec(
			reportId,
			m.MisconfigurationID,
			m.Type,
			m.Title,
			m.Description,
			m.Resolution,
			m.Severity,
			m.Target,
			m.Status,
		)
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

func (db *DB) GetVulnerabilities(reports []string, searchTerm string) ([]parser.Vulnerability, error) {
	size := len(reports)
	placeholders := make([]string, size)
	args := make([]interface{}, size+1)
	for i, r := range reports {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = r
	}

	args[size] = searchTerm
	searchTermPlaceholder := fmt.Sprintf("$%d", size+1)

	query := fmt.Sprintf(`
		SELECT id, vid, pkg_name, installed_version, fixed_version, title, description, severity, target 
													FROM vulnerability WHERE report_id IN (%s) 
											        		AND (%s = '' OR search_vector @@ plainto_tsquery(%s))
	ORDER BY id`, strings.Join(placeholders, ", "), searchTermPlaceholder, searchTermPlaceholder)

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vulnerabilities := make([]parser.Vulnerability, 0)
	for rows.Next() {
		var v parser.Vulnerability
		if err := rows.Scan(&v.ID, &v.VulnerabilityID, &v.PkgName, &v.InstalledVersion, &v.FixedVersion,
			&v.Title, &v.Description, &v.Severity, &v.Target); err != nil {
			return nil, err
		}

		vulnerabilities = append(vulnerabilities, v)
	}
	return vulnerabilities, nil
}

func (db *DB) GetMisconfigurations(reports []string, searchTerm string, resolved bool) ([]parser.Misconfiguration, error) {
	size := len(reports)
	placeholders := make([]string, size)
	args := make([]interface{}, 0, size+1)

	for i, r := range reports {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, r)
	}

	searchTermPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	args = append(args, searchTerm)

	query := fmt.Sprintf(`
		SELECT id, type, mid, title, description, resolution, severity, target, status
		FROM misconfiguration
		WHERE report_id IN (%s)
		  AND (%s = '' OR search_vector @@ plainto_tsquery(%s))`,
		strings.Join(placeholders, ", "),
		searchTermPlaceholder, searchTermPlaceholder,
	)

	if !resolved {
		query += " AND status IN ('FAIL', 'MANUAL')"
	}

	query += " ORDER BY id"

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	misconfigurations := make([]parser.Misconfiguration, 0)
	for rows.Next() {
		var m parser.Misconfiguration
		if err := rows.Scan(&m.ID, &m.Type, &m.MisconfigurationID, &m.Title, &m.Description, &m.Resolution, &m.Severity, &m.Target, &m.Status); err != nil {
			return nil, err
		}
		misconfigurations = append(misconfigurations, m)
	}
	return misconfigurations, nil
}

func (db *DB) DeleteMisconfigurationById(id string) error {
	query := `DELETE FROM misconfiguration WHERE id=$1`
	_, err := db.connection.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteMisconfigurationsBySearchTerm(reports []string, searchTerm string) error {
	size := len(reports)
	placeholders := make([]string, size)
	args := make([]interface{}, size+1)
	for i, r := range reports {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = r
	}

	args[size] = searchTerm
	searchTermPlaceholder := fmt.Sprintf("$%d", size+1)

	query := fmt.Sprintf(`
		DELETE FROM misconfiguration
		WHERE report_id IN (%s)
		  AND (%s = '' OR search_vector @@ plainto_tsquery(%s))
	`, strings.Join(placeholders, ", "), searchTermPlaceholder, searchTermPlaceholder)

	_, err := db.connection.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ResolveMisconfigurationById(id string) error {
	query := `
		UPDATE misconfiguration
		SET status = CASE original.status
		  WHEN 'FAIL' THEN 'RESOLVED'
		  WHEN 'MANUAL' THEN 'RESOLVED'
		  WHEN 'RESOLVED' THEN 'FAIL'
		  ELSE status
		END
		FROM (
		  SELECT mid, target, status
		  FROM misconfiguration
		  WHERE id = $1
		) AS original
		WHERE misconfiguration.mid = original.mid
		  AND misconfiguration.target = original.target;
`
	_, err := db.connection.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ResolveMisconfigurationsBySearchTerm(reportIds []string, searchTerm string) error {
	size := len(reportIds)
	placeholders := make([]string, size)
	args := make([]interface{}, size+1)
	for i, r := range reportIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = r
	}

	args[size] = searchTerm
	searchTermPlaceholder := fmt.Sprintf("$%d", size+1)

	query := fmt.Sprintf(`
		UPDATE misconfiguration m
		SET status = CASE
		  WHEN m.status = 'FAIL' THEN 'RESOLVED'
		  WHEN m.status = 'MANUAL' THEN 'RESOLVED'
		  WHEN m.status = 'RESOLVED' THEN 'FAIL'
		  ELSE m.status
		END
		FROM (
		  SELECT DISTINCT mid, target FROM misconfiguration
		  WHERE report_id IN (%s)
			AND (%s = '' OR search_vector @@ plainto_tsquery(%s))
			AND status IN ('FAIL', 'RESOLVED', 'MANUAL')
		) AS matched
		WHERE m.mid = matched.mid
		  AND m.target = matched.target
		  AND m.status IN ('FAIL', 'RESOLVED', 'MANUAL');
	`, strings.Join(placeholders, ", "), searchTermPlaceholder, searchTermPlaceholder)

	_, err := db.connection.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetReports(scanner string) ([]Report, error) {
	var reports []Report
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
		rows, err = db.connection.Query("SELECT id, filename FROM report where scanner_id=$1 order by generated_at desc", id)
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

func (db *DB) GetLastParsedReportIds() ([]string, error) {
	rows, err := db.connection.Query("SELECT DISTINCT ON (scanner_id) id FROM report ORDER BY scanner_id, generated_at DESC")
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	reportIds := make([]string, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		reportIds = append(reportIds, id)
	}

	return reportIds, nil
}

func (db *DB) GetScannerNameByReportId(reportId string) (string, error) {
	rows, err := db.connection.Query("SELECT scanner_id FROM report where id=$1", reportId)
	if err != nil {
		return "", fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("row scan failed: %v", err)
		}
		rows, err = db.connection.Query("SELECT name FROM scanner where id=$1", id)
		if err != nil {
			return "", fmt.Errorf("query failed: %v", err)
		}
		defer rows.Close()

		if rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				return "", fmt.Errorf("row scan failed: %v", err)
			}

			return name, nil
		}
	}

	return "", nil
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

package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

func writeToDatabase(db *sql.DB, address string, ports []string) (*CommittedScan, error) {
	intPorts, err := portsToIntSlice(ports)
	if err != nil {
		return nil, err
	}

	// We're rounding to the nearest second to keep things consistent with what
	// we get back from the database
	ucScan := UncommittedScan{address, time.Now().Round(time.Second), intPorts}
	id, err := insertScan(db, ucScan)
	if err != nil {
		return nil, err
	}

	scan := CommittedScan{id, ucScan.Host, ucScan.DateTime, ucScan.Ports}

	return &scan, nil
}

func insertScan(db *sql.DB, scan UncommittedScan) (int, error) {
	hostId, err := getHostId(db, scan.Host)
	if err != nil {
		return 0, err
	}

	ports := portsToStringSlice(scan.Ports)
	stmt, err := db.Prepare("INSERT INTO Scans (host_id, time, ports) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(hostId, scan.DateTime, strings.Join(ports, ","))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func getHostId(db *sql.DB, address string) (int, error) {
	// Does host exist in DB?
	row := db.QueryRow("SELECT id FROM Hosts WHERE address = ?", address)
	var hostId int = 0
	err := row.Scan(&hostId)
	if err != nil {
		log.Println(err.Error())
	}
	if hostId != 0 {
		return hostId, nil
	}

	// If not, create it
	_, err = db.Exec("INSERT INTO Hosts (address) VALUES (?)", address)
	if err != nil {
		// Error 1062 means the record already exists, which means another request
		// must have added it
		if err.(*mysql.MySQLError).Number != 1062 {
			return 0, err
		}
	}
	return getHostId(db, address)
}

func getScansForHost(db *sql.DB, address string, numScans int) ([]PartialCommittedScan, error) {
	// Determine which query to use
	sql := "SELECT Scans.id, time, ports FROM Scans JOIN Hosts ON Hosts.id = Scans.host_id WHERE address = ? ORDER BY Scans.id DESC"
	if numScans > 0 {
		sql = sql + " LIMIT " + strconv.Itoa(numScans)
	}

	// Run query
	results, err := db.Query(sql, address)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	// Convert results
	scans := []PartialCommittedScan{}
	for results.Next() {
		var id int = 0
		var timeString string
		var portListString string
		err = results.Scan(&id, &timeString, &portListString)
		if err != nil {
			return nil, err
		}

		stringPorts := strings.Split(portListString, ",")
		intPorts, err := portsToIntSlice(stringPorts)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse("2006-01-02 15:04:05", timeString)
		if err != nil {
			return nil, err
		}

		scans = append(scans, PartialCommittedScan{id, t, intPorts})
	}

	return scans, nil
}

func getPreviousScan(db *sql.DB, scanId int) (*CommittedScan, error) {
	sql := "SELECT Scans.id, time, ports, address FROM Scans JOIN Hosts ON Hosts.id = Scans.host_id JOIN (SELECT host_id FROM Scans WHERE id = ?) nested ON Hosts.id = nested.host_id WHERE Scans.id < ? ORDER BY Scans.id DESC LIMIT 1"
	scanIdString := strconv.Itoa(scanId)
	row := db.QueryRow(sql, scanIdString, scanIdString)

	var id int = 0
	var timeString string
	var portListString string
	var address string
	err := row.Scan(&id, &timeString, &portListString, &address)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return nil, err
	}

	stringPorts := strings.Split(portListString, ",")
	intPorts, err := portsToIntSlice(stringPorts)
	if err != nil {
		return nil, err
	}

	return &CommittedScan{id, address, t, intPorts}, nil
}

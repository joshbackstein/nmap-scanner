package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

func writeToDatabase(db *sql.DB, address string, ports []string) (*Scan, error) {
	intPorts, err := portsToIntSlice(ports)
	if err != nil {
		return nil, err
	}
	scan := Scan{address, time.Now(), intPorts}
	insertScan(db, scan)

	return &scan, nil
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

func insertScan(db *sql.DB, scan Scan) error {
	hostId, err := getHostId(db, scan.Host)
	if err != nil {
		return err
	}

	ports := portsToStringSlice(scan.Ports)
	_, err = db.Exec("INSERT INTO Scans (host_id, time, ports) VALUES (?,?,?)", hostId, scan.DateTime, strings.Join(ports, ","))
	if err != nil {
		return err
	}

	return nil
}

func getScansForHost(db *sql.DB, address string, numScans int) ([]PartialScan, error) {
	// Determine which query to use
	sql := "SELECT time, ports FROM Scans JOIN Hosts ON Hosts.id = Scans.host_id WHERE address = ? ORDER BY Scans.id DESC"
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
	scans := []PartialScan{}
	for results.Next() {
		var timeString string
		var portListString string
		err = results.Scan(&timeString, &portListString)
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

		scans = append(scans, PartialScan{t, intPorts})
	}

	return scans, nil
}

func getPreviousScanForHost(db *sql.DB, address string, dateTime time.Time) (*PartialScan, error) {
	sql := "SELECT time, ports FROM Scans JOIN Hosts ON Hosts.id = Scans.host_id WHERE address = ? AND time < ? ORDER BY Scans.id DESC LIMIT 1"
	row := db.QueryRow(sql, address, dateTime)

	var t time.Time
	var portListString string
	err := row.Scan(&t, &portListString)
	if err != nil {
		return nil, err
	}

	stringPorts := strings.Split(portListString, ",")
	intPorts, err := portsToIntSlice(stringPorts)
	if err != nil {
		return nil, err
	}

	return &PartialScan{t, intPorts}, nil
}

# Nmap Scanner

This Go project essentially creates a simple Nmap-as-a-Service that is backed by a MySQL database. It performs scans of ports 0-1000 on the hosts/IPs provided by the user and stores a history of those scans in the database. Additionally, the user is shown which ports are newly opened/closed between the current scan and the previous scan.

## Database

The table structure is located in `sql/setup-db.sql` and can be created like this:

```sh
# Replace "<user>" with your username
mysql --user=<user> --password < sql/setup-db.sql
```

## Building

These build instructions have only been tested on a Mac but should work on both Mac and Linux. To run the server, Go and Nmap will need to be installed. You can follow the instructions for your operating system to install them. You can build the server with:

```sh
cd src
go build -o nmap-scanner *.go
```

## Running

You might need to set some environment variables to get everything to run properly:

```sh
MYSQL_USERNAME=someone MYSQL_PASSWORD=password ./nmap-scanner
```

These are the available environment variables and their default values:
* HTTP\_PORT=8080
* MYSQL\_USERNAME=root
* MYSQL\_PASSWORD=password
* MYSQL\_HOST=localhost
* MYSQL\_PORT=3306
* MYSQL\_DB=NmapScanner

## Accessing the UI

After running the server, you can navigate to the UI in your web browser (default: `http://localhost:8080/`), enter a host/IP, and click "Scan"

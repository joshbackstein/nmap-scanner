CREATE DATABASE IF NOT EXISTS NmapScanner;

CREATE TABLE IF NOT EXISTS NmapScanner.Hosts
( id      INT           NOT NULL AUTO_INCREMENT
, address VARCHAR(253)  NOT NULL UNIQUE -- Max size of DNS hostname is 253 ASCII characters
, PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS NmapScanner.Scans
( id      INT           NOT NULL AUTO_INCREMENT
, host_id INT           NOT NULL
, time    DATETIME      NOT NULL
, ports   VARCHAR(4000) NOT NULL -- Comma-separated list of 0...1000 is 3894 characters long
, PRIMARY KEY (id)
, FOREIGN KEY (host_id) REFERENCES NmapScanner.Hosts(id)
);

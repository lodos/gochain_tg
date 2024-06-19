// common/clickhouse.go
package common

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
)

var ClickHouseConnect = map[string]string{
	"server":   "45.137.190.123",
	"user":     "dbu_license",
	"password": "Istanmul598.2007!",
	"db":       "license",
	"port":     "9000",
}

type ClickHouseConnection struct {
	db *sql.DB
}

//func NewClickHouseConnection() *ClickHouseConnection {
//	return &ClickHouseConnection{
//		db: nil, // Initialize it to nil
//	}
//}

func (c *ClickHouseConnection) Init() error {
	// Open a connection to ClickHouse
	connect, err := sql.Open("clickhouse", fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20",
		ClickHouseConnect["server"], ClickHouseConnect["port"], ClickHouseConnect["user"],
		ClickHouseConnect["password"], ClickHouseConnect["db"],
	))
	if err != nil {
		return err
	}

	c.db = connect
	return nil
}

func (c *ClickHouseConnection) Query(query string) (*sql.Rows, error) {
	// Execute the query using the ClickHouse connection
	return c.db.Query(query)
}

func (c *ClickHouseConnection) Close() {
	// Close the ClickHouse connection
	if c.db != nil {
		c.db.Close()
	}
}

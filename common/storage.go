package common

import (
	"database/sql"
	"fmt"
)

type EventStorage struct {
	ClickHouseConnect map[string]string
}

// Init initializes the EventStorage with ClickHouse connection details
func (es *EventStorage) Init(server, user, password, db string, port int) {
	es.ClickHouseConnect = map[string]string{
		"server":   server,
		"user":     user,
		"password": password,
		"db":       db,
		"port":     fmt.Sprintf("%d", port),
	}
}

// InsertToClickHouse inserts the given query into ClickHouse
func (es *EventStorage) GetJobsForTask(query string) error {
	connect, err := sql.Open("clickhouse", fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20",
		es.ClickHouseConnect["server"], es.ClickHouseConnect["port"], es.ClickHouseConnect["user"],
		es.ClickHouseConnect["password"], es.ClickHouseConnect["db"],
	))
	if err != nil {
		return err
	}
	defer func(connect *sql.DB) {
		err := connect.Close()
		if err != nil {

		}
	}(connect)

	_, err = connect.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

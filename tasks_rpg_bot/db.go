package main

import (
	_ "github.com/mattn/go-sqlite3"
	//"github.com/mattn/go-sqlite3"
	"database/sql"
)

type DbManager struct {
	dsn string // DB connection settings
	//db *sqlite3.SQLiteConn // Instance of connection
	db *sql.DB // Instance of connection
}

func (m *DbManager) init(dsn string) {
	m.dsn = dsn
	//m.db = sqlite3.con

	conn, err := sql.Open(`sqlite3`, dsn)

	if err != nil {
		logger.Fatal(`Failed to connect to "%s": %s`, dsn, err)
	}

	m.db = conn

	//logger.Fatal(`Connection %v`, *db)

	m.initTables()
}

func (m *DbManager) initTables() {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY, user_id TEXT, title TEXT, description TEXT, status TEXT, exp INTEGER, date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)",
	}

	for _, query := range queries {
		logger.Debug(`Running SQL: %s`, query)

		statement, err := m.db.Prepare(query)
		if err != nil {
			logger.Error("Failed executing SQL query: \n%s\nSQL Error: %s", query, err)
		}
		statement.Exec()
	}
}

func NewDbManager(dsn string) *DbManager {
	var manager DbManager

	manager.init(dsn)

	return &manager
}

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
		"CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, title TEXT, description TEXT, status TEXT, exp INTEGER, date_expiration TEXT DEFAULT '', date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)",
		"CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, name TEXT, exp INTEGER, date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)",
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

func (m *DbManager) findFutureTasks() (entities []TaskDbEntity) {
	rows, err := m.db.Query("SELECT * FROM task WHERE date_expiration > datetime('now') and status = $1", statusPending)

	if err != nil {
		logger.Error(`SQL error: %s`, err)
	}

	for rows.Next() {
		var taskEntity TaskDbEntity
		err = taskEntity.Load(rows)

		if err != nil {
			logger.Info(`SQL error data load: %s`, err)

			continue
		}

		entities = append(entities, taskEntity)
	}

	return entities
}

func (m *DbManager) findAllTasks() (entities []TaskDbEntity) {
	rows, err := m.db.Query("SELECT * FROM task")

	if err != nil {
		logger.Error(`SQL error: %s`, err)
	}

	for rows.Next() {
		var taskEntity TaskDbEntity
		err = taskEntity.Load(rows)

		if err != nil {
			logger.Info(`SQL error data load: %s`, err)

			continue
		}

		//if logger.DebugLevel() {
		//	logger.Debug(`SQL ROW: %s`, encodeToJson(taskEntity))
		//}

		entities = append(entities, taskEntity)
	}

	return entities
}

func (m *DbManager) findTasksByUser(userId int) (entities []TaskDbEntity) {
	rows, err := m.db.Query("SELECT * FROM task WHERE user_id = $1", userId)
	if err != nil {
		logger.Error(`SQL error: %s`, err)
	}

	for rows.Next() {
		var taskEntity TaskDbEntity
		err = taskEntity.Load(rows)

		if err != nil {
			logger.Info(`SQL error data load: %s`, err)

			continue
		}

		entities = append(entities, taskEntity)
	}

	return entities
}

func (m *DbManager) findTaskById(id int) *TaskDbEntity {
	rows, err := m.db.Query("SELECT * FROM task WHERE id = $1", id)
	if err != nil {
		logger.Error(`SQL error: %s`, err)
	}

	for rows.Next() {
		var taskEntity TaskDbEntity
		err = taskEntity.Load(rows)

		if err != nil {
			logger.Info(`SQL error data load: %s`, err)

			return nil
		}

		return &taskEntity
	}

	return nil
}

func NewDbManager(dsn string) *DbManager {
	var manager DbManager

	manager.init(dsn)

	return &manager
}

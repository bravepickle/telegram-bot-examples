package main

import (
	"database/sql"
	"time"
)

const statusPending = `pending`
const statusDone = `done`
const statusCanceled = `canceled`

type DbTime struct{ time.Time }

func (t DbTime) String() string {
	return t.Format(`2006-01-02 15:04:05`)
}

func NewDbTime(t time.Time) DbTime {
	// change timezone to local time to have single tz everywhere
	// maximum precision - second
	return DbTime{Time: t.Local().Round(time.Second)}
}

// list of all available entities for application

type DbEntityInterface interface {
	Save() bool
	Load(sqlRows *sql.Rows) error
	//Init()
	Table() string
}

type DbEntityStruct struct {
	Id int
	//InitSql string // string to init table for given entity
}

func (e DbEntityStruct) isNewRecord() bool {
	return e.Id == 0
}

func (e DbEntityStruct) Table() string {
	return `[undefined]`
}

type TaskDbEntity struct {
	Id             int
	UserId         int
	Title          string
	Status         string
	Exp            int
	Description    string
	DateExpiration DbTime
	DateCreated    DbTime
	DateUpdated    DbTime

	DbEntityStruct
	//DbEntityInterface
}

func (e *TaskDbEntity) Load(sqlRows *sql.Rows) error {
	// TODO: list fields to change params list
	return sqlRows.Scan(&e.Id, &e.UserId, &e.Title, &e.Description, &e.Status, &e.Exp, &e.DateExpiration, &e.DateCreated, &e.DateUpdated)
	//return sqlRows.Scan(&e.Id)
}

func (e *TaskDbEntity) Save() bool {
	if !e.isNewRecord() {
		logger.Fatal(`Not implemented update method for entity "%s"`, e.Table())
	}

	//database.Prepare("INSERT INTO task (user_id, title, status, description, status, exp) VALUES (?, ?, ?, ?, ?, ?)")
	statement, err := dbManager.db.Prepare(
		"INSERT INTO task (user_id, title, status, description, exp, date_expiration, date_created, date_updated) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		logger.Error(`SQL error: %s`, err)
	}

	result, err := statement.Exec(e.UserId, e.Title, e.Status, e.Description, e.Exp, e.DateExpiration.String(), e.DateCreated.String(), e.DateUpdated.String())
	if err != nil {
		logger.Error(`SQL error: %s`, err)

		return false
	}

	// TODO: other values update as well (refresh entity)
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		logger.Error(`SQL last insert ID read error: %s`, err)

		return false
	}

	e.Id = int(lastInsertId)
	if logger.DebugLevel() {
		logger.Debug(`SQL result: %s`, encodeToJson(e))
	}

	return true
}

func (e TaskDbEntity) Table() string {
	return `task`
}

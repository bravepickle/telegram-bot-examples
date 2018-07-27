package main

import "database/sql"

// list of all available entities for application

type DbEntityInterface interface {
	Save() bool
	Load(sqlRows *sql.Rows) error
	Init()
}

type DbEntityStruct struct {
	Id int
	//InitSql string // string to init table for given entity
}

func (e DbEntityStruct) isNewRecord() bool {
	return e.Id == 0
}

type TaskDbEntity struct {
	Id             int
	UserId         int
	Title          string
	Status         string
	Exp            int
	Description    string
	DateExpiration string
	DateCreated    string
	DateUpdated    string

	DbEntityStruct
	DbEntityInterface
}

//func (e *TaskDbEntity) Init() {
//	e.InitSql = "CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY, user_id TEXT, title TEXT, description TEXT, status TEXT, exp INTEGER, date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)"
//}

//
//func (e *TaskDbEntity) Save() {
//
//	dbManager.db.
//}

func (e *TaskDbEntity) Load(sqlRows *sql.Rows) error {
	// TODO: list fields to change params list
	return sqlRows.Scan(&e.Id, &e.UserId, &e.Title, &e.Description, &e.Status, &e.Exp, &e.DateExpiration, &e.DateCreated, &e.DateUpdated)
	//return sqlRows.Scan(&e.Id)
}

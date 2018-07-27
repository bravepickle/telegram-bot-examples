package main

import "database/sql"

// list of all available entities for application

type DbEntityInterface interface {
	Save() bool
	Scan(sqlRows *sql.Rows) error
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
	DbEntityStruct
	DbEntityInterface

	Id          int
	Title       string
	Status      string
	Description string
	DateCreated string
	DateUpdated string
}

//func (e *TaskDbEntity) Init() {
//	e.InitSql = "CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY, user_id TEXT, title TEXT, description TEXT, status TEXT, exp INTEGER, date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)"
//}

//
//func (e *TaskDbEntity) Save() {
//
//	dbManager.db.
//}

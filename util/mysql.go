package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//go get github.com/go-sql-driver/mysql 

var db *sql.DB 
var err error

func init () {
	db, err = sql.Open("mysql", "dbuser:dbpasswd@tcp(dburl:3306)/cm_launcher_theme?charset=utf8")  
    if err != nil {
        log.Fatal("sql.Open ", err)  
    } 
}

func QueryOneRow(sqlstr string, values ...interface{}) error {
	//var row *sql.Row
    row := db.QueryRow(sqlstr)
    err = row.Scan(values...)
	return err
}
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		log.Fatal(err)
	}

	// 用户表
	_, _ = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`)

	// 上传记录表
	_, _ = DB.Exec(`
	CREATE TABLE IF NOT EXISTS uploads(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		filepath TEXT,
		user_id INTEGER,
		author TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
}

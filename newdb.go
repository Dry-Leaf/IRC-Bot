package main

import (
    "database/sql"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

const (
    createMessageTabelSQL = `CREATE TABLE messages (
        "Recipient" TEXT NOT NULL,
        "Sender" TEXT NOT NULL,
        "Message" TEXT NOT NULL
    );`

    createWetTabelSQL = `CREATE TABLE wet (
        "User" TEXT PRIMARY KEY,
        "Location" TEXT NOT NULL
    );`

    createHostTableSQL = `CREATE TABLE hosts (
        "Host" TEXT PRIMARY KEY
    );`

    createUserTableSQL = `CREATE TABLE users (
        "Username" TEXT PRIMARY KEY
    );`
)

func create_table(db *sql.DB) {

    statement, err := db.Prepare(createMessageTabelSQL)
    Err_check(err)
    statement.Exec()

    statement, err = db.Prepare(createWetTabelSQL)
    Err_check(err)
    statement.Exec()
    
    statement, err = db.Prepare(createHostTableSQL)
    Err_check(err)
    statement.Exec()

    statement, err = db.Prepare(createUserTableSQL)
    Err_check(err)
    statement.Exec()
}

func New_db() {

    file, err := os.Create("message.db")
    Err_check(err)

    file.Close()

    conn, err := sql.Open("sqlite3", "message.db")
    Err_check(err)
    defer conn.Close()

    create_table(conn)
}

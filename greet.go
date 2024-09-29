package main

import (
    "database/sql"

    "github.com/thoj/go-ircevent"
    _ "github.com/mattn/go-sqlite3"        
)

const (
    host_query = `SELECT EXISTS (SELECT 1 FROM hosts where Host = ?)`
    insert_host = `INSERT INTO hosts (Host) VALUES (?)`

    greeting_1 = `Welcome to the `
    greeting_1b = ` IRC channel desu.`
)

func Greet(host, ch string, conn *irc.Connection) {
    dbconn, err := sql.Open("sqlite3", "message.db")
    Err_check(err)
    defer dbconn.Close()

    tx, _ := dbconn.Begin()

    var repeat_user bool
    err = tx.QueryRow(host_query, host).Scan(&repeat_user)
    Err_check(err)

    if !repeat_user {
        tx.Exec(insert_host, host)

        conn.Privmsg(ch, greeting_1 + ch + greeting_1b)
    }
    
    tx.Commit()
}

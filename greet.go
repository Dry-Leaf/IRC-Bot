package main

import (
    "time"
    "database/sql"

    "github.com/thoj/go-ircevent"
    _ "github.com/mattn/go-sqlite3"        
)

const (
    host_query = `SELECT EXISTS (SELECT 1 FROM hosts where Host = ?)`
    insert_host = `INSERT INTO hosts (Host) VALUES (?)`

    greeting_1 = `Welcome to Tohno-chan's official IRC channel.`
    greeting_2 = `Feel free to chat about anything within reason.`
    greeting_3 = `You can use /nick to change your nickname from the default, brohno.`
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

        conn.Privmsg(ch, greeting_1)
        time.Sleep(time.Second)
        conn.Privmsg(ch, greeting_2)
        time.Sleep(time.Second)
        conn.Privmsg(ch, greeting_3)
    }
    
    tx.Commit()
}

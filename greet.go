package main

import (
    "strings"
    "database/sql"

    "github.com/thoj/go-ircevent"
    _ "github.com/mattn/go-sqlite3"        
)

const (
    user_query = `SELECT EXISTS (SELECT 1 FROM users where Username = ?)`
    host_query = `SELECT EXISTS (SELECT 1 FROM hosts where Host = ?)`

    insert_user = `INSERT INTO users (Username) VALUES (?)`
    insert_host = `INSERT INTO hosts (Host) VALUES (?)`

    greeting_1 = `Welcome to the `
    greeting_1b = ` IRC channel desu.`
)

func Greet(username_full, host, ch string, conn *irc.Connection) {
    username := strings.TrimRight(username_full, "_")
    
    _, is_default := Def_usernames[username]
    if is_default {
        conn.Privmsg(ch, greeting_1 + ch + greeting_1b)
        return
    }

    dbconn, err := sql.Open("sqlite3", "message.db")
    Err_check(err)
    defer dbconn.Close()

    tx, _ := dbconn.Begin()

    var repeat_user bool
    err = tx.QueryRow(user_query, username).Scan(&repeat_user)
    Err_check(err)

    if repeat_user {
        return
    } else {
        tx.Exec(insert_user, username)
    }

    err = tx.QueryRow(host_query, host).Scan(&repeat_user)
    Err_check(err)

    if !repeat_user {
        tx.Exec(insert_host, host)

        conn.Privmsg(ch, greeting_1 + ch + greeting_1b)
    }
    
    tx.Commit()
}

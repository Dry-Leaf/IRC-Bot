package main 

import (
    "database/sql"
    //"fmt"

    "regexp"
    "github.com/thoj/go-ircevent"
    _ "github.com/mattn/go-sqlite3"
)

var newmessage_reg = regexp.MustCompile(`(?i)\A\.tell(?:\s+)(\S+)(?:\s+)(.+)`) 

func NewMessage(sender string, input string, conn *irc.Connection) {
    params := newmessage_reg.FindStringSubmatch(input)

    if len(params) == 0 {return}

    dbconn, err := sql.Open("sqlite3", "message.db")
    Err_check(err)
    defer dbconn.Close()

    statement, err := dbconn.Prepare(`INSERT INTO messages (Sender, Recipient, Message) 
        VALUES (?, ?, ?)`)
    Err_check(err)
    statement.Exec(sender, params[1], params[2])
}

func Tell(user string, conn *irc.Connection) {
    dbconn, err := sql.Open("sqlite3", "message.db")
    Err_check(err)
    defer dbconn.Close()

    tx, _ := dbconn.Begin()
    statement := `SELECT Sender, Message FROM messages where Recipient = ? COLLATE NOCASE`

    message_rows, err := tx.Query(statement, user)
    if err != nil {tx.Rollback()} 

    defer message_rows.Close()

    if err == nil {
        for message_rows.Next() {
            var sender string
            var message string

            err = message_rows.Scan(&sender, &message)
            Err_check(err)
            conn.Privmsg(user, sender + ": " + message)

            tx.Exec("DELETE FROM messages WHERE Recipient = ? COLLATE NOCASE AND Sender = ? COLLATE NOCASE AND Message = ?", user, sender, message)
    }}

    tx.Commit()
}
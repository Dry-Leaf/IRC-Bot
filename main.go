package main
//An IRC bot written for tohno-chan

import (
    "log"
    "os"
    "time"
    "strings"
    //"fmt"

    "github.com/thoj/go-ircevent"
)


func Err_check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

//to circumvent word filters
func Vowel_replace(input string) string {
    //cyrillic
    output := strings.ReplaceAll(input, "a", "а")
    output = strings.ReplaceAll(output, "A", "А")

    output = strings.ReplaceAll(output, "e", "е")
    output = strings.ReplaceAll(output, "E", "Е")

    output = strings.ReplaceAll(output, "i", "і")
    output = strings.ReplaceAll(output, "I", "І")

    output = strings.ReplaceAll(output, "o", "о")
    output = strings.ReplaceAll(output, "O", "О")

    return output
}

func main() {

    log_file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    Err_check(err)
    defer log_file.Close()

    log.SetOutput(log_file)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    Load_conf()
    Trivia_unload()

    if _, err = os.Stat("message.db"); err != nil {
        New_db()
    }

    conn := irc.IRC(Nickname, Username)
    conn.Log.SetOutput(log_file)
    //conn.Debug = true

    //connecting to the server
    channel_list := strings.Join(Channels, ",")
    conn.AddCallback("001", func(e *irc.Event) {
            conn.Join(channel_list)
    })

    //anybody posts in the channel
    conn.AddCallback("PRIVMSG", func (e *irc.Event) {
        stored := e.Message()

        go Youtube(stored, e.Arguments[0], conn)
        go Openweather(stored, e.Arguments[0], conn)
        go Urbandictionary(stored, e.Arguments[0], conn)
        go EightBall(stored, e.Arguments[0], conn)
        go NewMessage(e.Nick, stored, conn)
        go Trivia(e.Nick, stored, e.Arguments[0], conn)
    })

    //anybody joins the server
    conn.AddCallback("JOIN", func (e *irc.Event) {
        user := e.Nick
        go Tell(user, conn)
    })

    //kicked
    conn.AddCallback("KICK", func (e *irc.Event) {
        kick_args := e.Arguments

        if kick_args[1] == conn.GetNick() {
            time.Sleep(3 * time.Second)
            conn.Join(kick_args[0])
        }
    })

    err = conn.Connect(Server)
    Err_check(err)
    defer conn.Quit()

    conn.Loop()
}

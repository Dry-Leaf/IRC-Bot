package main

import (
    "math/rand"
    "regexp"
    
    "github.com/thoj/go-ircevent"
)

var eightBallAnswers = map[int]string{
    0: "It is certain.", 1: "It is decidedly so.",
    2: "Yes definitely.", 3: "You may rely on it.",
    4: "As I see it, yes.", 5: "Most likely.",
    6: "Outlook good.", 7: "Yes.",
    8: "Signs point to yes.", 9: "Reply hazy, try again.",
    10: "Ask again later.", 11: "Better not tell you now.",
    12: "Cannot predict now.", 13: "Concentrate and ask again.",
    14: "Don't count on it.", 15: "My reply is no.",
    16: "My sources say no.", 17: "Outlook not so good.",
    18: "Very doubtful.", 19: "Without a doubt."}

var eightBallReg = regexp.MustCompile(`(?i)\A\.8ball(?:\s+|\z)`)

func EightBall(stored, ch string, conn *irc.Connection) {
    if eightBallReg.MatchString(stored) {
        reply := eightBallAnswers[rand.Intn(20)]
        conn.Privmsg(ch, reply)
    }
}

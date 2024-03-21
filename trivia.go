package main

import (
    "math/rand"
    "html"
    //"os"
    "strconv"
    "slices"
    "strings"
    "fmt"
    "regexp"
    "time"
    "encoding/json"
    "io/ioutil"

    "github.com/thoj/go-ircevent"
)

var playing_trivia = false

var triviaReg = regexp.MustCompile(`(?i)\A\.triv(?:\s+|\z)(\d*)`)
var striviaReg = regexp.MustCompile(`(?i)\A\.striv(?:\s+|\z)`)

type question_struct struct {
    Question string       `json:"question"`
    Answers []string      `json:"answers"`
}

var qslice []question_struct
var selection []int
var score = map[string]int{}

var reply = make(chan string)
var asking = make(chan bool)
var hinting = make(chan bool)
var answered = make(chan bool)

func Trivia_unload() {
    file, err := ioutil.ReadFile("interleaved.json")
    Err_check(err)

    err = json.Unmarshal(file, &qslice)
    Err_check(err)
}

func reset() {
    playing_trivia = false
    selection = nil
    clear(score)
}

func striv() {
    reset()
    asking <- false
    hinting <- false
}

func check(sender, ch string, conn *irc.Connection) {
    submission := <-reply
    answers := qslice[selection[0]].Answers

    for _, answer := range answers {
        if strings.EqualFold(answer, submission) {
            asking <- false
            hinting <- false
            conn.Privmsg(ch, fmt.Sprintf("Winner: %s | Answer: %s", sender, answers[0]))

            if _, ok := score[sender]; ok {
                score[sender] += 1
            } else {
                score[sender] = 1
            }
            return
    }}
}

func declare_winner(ch string, conn *irc.Connection) {
    var winner string
    points := 0
    for k, v := range score {
        if v > points {
            winner = k; points = v
        }
    }

    conn.Privmsg(ch, fmt.Sprintf("Trivia Complete. Winner: %s | Points: %d", winner, points))
}

func ask(ch string, conn *irc.Connection) {
    hint_give := func(hint, answer []rune) {

        hint_size := len(hint)
        timer := time.After(12 * time.Second)

        select {
            case <-timer:
                if string(hint) == string(answer) {
                    conn.Privmsg(ch, fmt.Sprintf("Times Up! The answer is: %s", string(hint)))
                    asking <- false
                    return
                }

                conn.Privmsg(ch, fmt.Sprintf("Hint: %s", string(hint)))
                for {
                    num := rand.Intn(hint_size)
                    if hint[num] == '*' {hint[num] = answer[num]; break}
                }

                for i, _ := range hint {
                    randomNumber := rand.Float64()
                    if randomNumber < .4 {hint[i] = answer[i]}
                }
                asking <- true
            case <-hinting:
    }}

    question_loop := func(cq question_struct, hint []rune) {
        begin := <-asking

        for begin {
            go hint_give(hint, []rune(cq.Answers[0]))
            begin = <-asking
        }

        if len(selection) > 0 {
            selection = selection[1:]
            time.Sleep(3 * time.Second)
        }
        answered <- true
    }

    for qc := 1; playing_trivia; qc++ {
        if len(selection) < 1 {
           declare_winner(ch, conn)
           reset()
           break
        }

        cq := qslice[selection[0]]
        conn.Privmsg(ch, fmt.Sprintf("%d. %s", qc, html.UnescapeString(cq.Question)))

        var hint []rune
        for _, c := range cq.Answers[0] {
            if c == ' ' {
                hint = append(hint, ' ')
            } else {hint = append(hint, '*')}
        }
        
        go question_loop(cq, hint)
        asking <- true
        <-answered
    }
}

func Trivia(sender, stored, ch string, conn *irc.Connection) {    
    if playing_trivia && striviaReg.MatchString(stored) {
        declare_winner(ch, conn)
        go striv()
        return
    }
    
    params := triviaReg.FindStringSubmatch(stored)

    if playing_trivia {
        go check(sender, ch, conn)
        reply <- stored
    } else if triviaReg.MatchString(stored) {      
        question_number := 0
        
        if params[1] != "" {
            question_number, _ = strconv.Atoi(params[1])
        } else {
            question_number = 5
        }
        
        pool_size := len(qslice)

        if (question_number <= 0) {
            conn.Privmsg(ch, "Question number too low")
            return
        }
        
        if (question_number > pool_size) {
            conn.Privmsg(ch, "Question number must not exceed: " + string(pool_size))
            return
        }

        for len(selection) < question_number {
            num := rand.Intn(pool_size)
            if !slices.Contains(selection, num) {selection = append(selection, num)}
        }

        playing_trivia = true
        go ask(ch, conn)
    }
}


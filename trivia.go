package main

import (
    "math/rand"
    //"html"
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
//var score = map[string]int{}

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
    //clear(score)
    select {
        case <-asking:
        default:
    }
    select {
        case <-hinting:
        default:
    }
    select {
        case <-answered:
        default:
    }
    select {
        case <- reply:
        default:
    }
}

func check(sender, ch string, conn *irc.Connection) {
    submission := <-reply
    answers := qslice[selection[0]].Answers

    for _, answer := range answers {
        if strings.EqualFold(answer, submission) {
            asking <- false
            hinting <- false
            conn.Privmsg(ch, fmt.Sprintf("Winner %s | Answer: %s", sender, answers[0]))
            return
    }}
}

func ask(ch string, conn *irc.Connection) {
    hint_give := func(hint, answer []rune) {
        hint_size := len(hint)
        timer := time.After(15 * time.Second)

        select {
            case <-timer:
                conn.Privmsg(ch, fmt.Sprintf("Hint: %s", string(hint)))
                for {
                    num := rand.Intn(hint_size)
                    if hint[num] == '*' {hint[num] = answer[num]; break}
                }
                asking <- true
            case <-hinting:
    }}

    question_loop := func(cq question_struct, hint []rune) {
        begin := <-asking

        for begin {
            if string(hint) == cq.Answers[0] {
                conn.Privmsg(ch, fmt.Sprintf("Times Up! The answer is: %s", string(hint)))
                break
            }
            go hint_give(hint, []rune(cq.Answers[0]))
            begin = <-asking
        }

        selection = selection[1:]
        answered <- true
    }

    for playing_trivia {
        if len(selection) < 1 {
           reset()
           return
        }

        cq := qslice[selection[0]]
        conn.Privmsg(ch, cq.Question)

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
    /*if striviaReg.MatchString(stored) {
        conn.Privmsg(ch, "Trivia Quit")
        reset()
        return
    }*/
    
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


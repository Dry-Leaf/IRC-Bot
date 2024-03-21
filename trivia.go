package main

import (
    "math"
    "math/rand"
    "html"
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
    select {
        case asking <- false:
        default:
    } 
    select {
        case hinting <- false:
        default:
}}

func check(sender, submission, ch string, conn *irc.Connection) {
    if len(selection) < 1 {return}
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
    var winners []string
    points := 0
    for k, v := range score {
        if v >= points {
            if v > points {winners = nil; points = v}
            winners = append(winners, k)
    }}

    message := "Trivia Complete. Winner: %s | Points: %d"
    if len(winners) > 1 {message = "Trivia Complete. Tie Between: %s | Points: %d"}

    conn.Privmsg(ch, fmt.Sprintf(message, strings.Join(winners, ", "), points))
}

func ask(ch string, conn *irc.Connection) {
    hint_give := func(hint, answer []rune, iter int) {

        hint_size := len(hint)
        timer := time.After(12 * time.Second)

        select {
            case <-timer:
                if iter == 3 || (iter > 0 && 0 >= hint_size - 3) {
                    asking <- false
                    conn.Privmsg(ch, fmt.Sprintf("Times Up! The answer is: %s", string(answer)))
                    return
                }
                
                if hint_size > 3 {
                    ps := int(max(math.Floor(float64(hint_size) / 4), 1))

                    switch iter {
                    case 0:
                        copy(hint[0:], answer[0:ps])
                    case 1:
                        copy(hint[ps:], answer[ps:2*ps])
                    case 2:
                        copy(hint[ps*3:], answer[3*ps:4*ps])
                }}            

                conn.Privmsg(ch, fmt.Sprintf("Hint: %s", string(hint)))
                asking <- true
            case <-hinting:
    }}

    question_loop := func(cq question_struct, hint []rune) {
        begin := <-asking

        for iter := 0; begin; iter++ {
            go hint_give(hint, []rune(cq.Answers[0]), iter)
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
            if c == ' ' || c == ',' || c == '\'' || c == '%' || c == '-'{
                hint = append(hint, c)
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
        go check(sender, stored, ch, conn)
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


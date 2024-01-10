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
    "encoding/json"
    "io/ioutil"

    "github.com/thoj/go-ircevent"
)

var playing_trivia = false
var asking = false
var answer = "42"

var triviaReg = regexp.MustCompile(`(?i)\A\.triv(?:\s+|\z)(\d*)`)
var striviaReg = regexp.MustCompile(`(?i)\A\.striv(?:\s+|\z)`)
var lett_arr = [4]string{"a", "b", "c", "d"}

var score = map[string]int{}

type question_struct struct {
    Question string        `json:"question"`
    Answer string          `json:"answer"`
    Incorrect []string     `json:"incorrect"`
    Qtype string           `json:"type"`
}

var qslice []question_struct
var selection []int

func Trivia_unload() {
    file, err := ioutil.ReadFile("interleaved.json")
    Err_check(err)

    err = json.Unmarshal(file, &qslice)
    Err_check(err)
}

func ask(ch string, conn *irc.Connection) {
    cq := qslice[selection[0]]

    conn.Privmsg(ch, html.UnescapeString(cq.Question))

    if cq.Qtype == "boolean" {
        conn.Privmsg(ch, "True or False?")
        answer = strings.ToLower(cq.Answer[0:1])
    } else { 
        var possible []string
        possible = append(cq.Incorrect, cq.Answer)
        rand.Shuffle(len(possible), func(i, j int) { possible[i], possible[j] = possible[j], possible[i] })

        for i, e := range possible {
           conn.Privmsg(ch, fmt.Sprintf("%s: %s", lett_arr[i], e))
           if cq.Answer == possible[i] {answer = lett_arr[i]}
    }}  
}


func Trivia(sender, stored, ch string, conn *irc.Connection) {
    reset := func() {
        playing_trivia = false
        asking = false
        clear(score)
        clear(selection)
    }
    
    params := triviaReg.FindStringSubmatch(stored)
    
    if triviaReg.MatchString(stored) && !playing_trivia {      
        question_number := 0
        
        playing_trivia = true

        if params[1] != "" {
            question_number, _ = strconv.Atoi(params[1])
        } else {
            question_number = 10
        }
        
        pool_size := len(qslice)
        fmt.Println(pool_size)
        
        for len(selection) < question_number {
            num := rand.Intn(pool_size)
            if !slices.Contains(selection, num) {selection = append(selection, num)}
        }
    }

    if striviaReg.MatchString(stored) {reset()}

    if playing_trivia {
        if len(selection) < 1 {
           reset()
           return
        }
       
        if !asking {
            asking = true
            ask(ch, conn)
        } else {
            correct := false
            
            if strings.ToLower(stored[0:1]) == answer {correct = true}

            if correct {
                asking = false
                conn.Privmsg(ch, fmt.Sprintf("Correct! The answer is %s", answer))

                selection = selection[1:]
                if _, ok := score[sender]; ok {
                    score[sender] += 1
                } else {
                    score[sender] = 1
                }
                ask(ch, conn)
        }}
}}

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

var triviaReg = regexp.MustCompile(`(?i)\A\.triv(?:\s+|\z)(\d*)`)
var striviaReg = regexp.MustCompile(`(?i)\A\.striv(?:\s+|\z)`)

var score = map[string]int{}

type question_struct struct {
    Question string        `json:"question"`
    Answer string          `json:"answer"`
    Incorrect []string     `json:"incorrect"`
    Qtype string           `json:"type"`
}

var qslice []question_struct

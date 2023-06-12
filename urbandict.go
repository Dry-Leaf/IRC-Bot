package main
//youtube feature

import (
    "regexp"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strings"
    "github.com/thoj/go-ircevent"
)

type UDefinition struct {
    Definition string     `json:"definition"`
}

var urbandict_reg = regexp.MustCompile(`(?i)\A\.ud(?:\s+|\z)(\S.*)`)

func Urbandictionary(stored, ch string, conn *irc.Connection) {
    tolookup := urbandict_reg.FindStringSubmatch(stored)

    if len(tolookup) > 0{
        api_url := `https://api.urbandictionary.com/v0/define?term=` + strings.ReplaceAll(tolookup[1], " ", "%20")
        var client = &http.Client{Timeout: 10 * time.Second}

        resp, err := client.Get(api_url)
        if err != nil {
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {            
            body, err := ioutil.ReadAll(resp.Body)
            Err_check(err)
            
            var dat map[string]json.RawMessage
            err = json.Unmarshal(body, &dat)
            Err_check(err)

            var ud []UDefinition
            err = json.Unmarshal(dat["list"], &ud)
            Err_check(err)

            if len(ud) > 0 {
                conn.Privmsg(ch, ud[0].Definition)
            } else {
                conn.Privmsg(ch, "Word not found.")
            }
        }
    }
}
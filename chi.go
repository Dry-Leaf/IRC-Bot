package main
//An IRC bot written for tohno-chan

import (
    "log"
    "regexp"
    "os"
    "net/http"
    "encoding/json"
    "time"
    "io/ioutil"

    "github.com/thoj/go-ircevent"
)


func Err_check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

var youtube_linkreg = regexp.MustCompile(`(https:\/\/|https:\/\/www.)(youtube.com\/watch\?v=|youtu.be\/)(\S+)`)

//posts title and description of youtube videos
func youtube(stored string, conn *irc.Connection) {
    url := youtube_linkreg.FindStringSubmatch(stored)

    if len(url) == 4 {
        id := url[3]

        var client = &http.Client{Timeout: 10 * time.Second}

        api_url := "https://www.googleapis.com/youtube/v3/videos?part=snippet&fields=items/snippet/description,items/snippet/title&id=" +
                    id + "&key=" + YT_apikey
        resp, err := client.Get(api_url)
        Err_check(err)
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            var dat map[string]interface{}
            body, err := ioutil.ReadAll(resp.Body)
            Err_check(err)
            err = json.Unmarshal(body, &dat)
            Err_check(err)

            items := dat["items"].([]interface{})[0]
            snippets := items.(map[string]interface{})["snippet"]
            title := snippets.(map[string]interface{})["title"].(string)
            desc := snippets.(map[string]interface{})["description"].(string)

            conn.Privmsg(Channel, title + " ★ " + desc)
        }
    }
}


func main() {

    log_file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    Err_check(err)
    defer log_file.Close()

    log.SetOutput(log_file)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    Load_conf()

    conn := irc.IRC(Nickname, Username)
    conn.Log.SetOutput(log_file)
    //conn.Debug = true

    //connecting to the server
    conn.AddCallback("001", func(e *irc.Event) {
        conn.Join(Channel) 
    })

    //anybody posts in the channel
    conn.AddCallback("PRIVMSG", func (e *irc.Event) {
        stored := e.Message()

        go youtube(stored, conn)
    })

    //kicked
    conn.AddCallback("KICK", func (e *irc.Event) {
        kick_args := e.Arguments

        if kick_args[1] == conn.GetNick() {
            time.Sleep(3 * time.Second)
            conn.Join(Channel)
        }
    })

    err = conn.Connect(Server)
    Err_check(err)
    defer conn.Quit()

    conn.Loop()
}

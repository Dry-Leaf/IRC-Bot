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
    "fmt"

    "github.com/thoj/go-ircevent"
)


func Err_check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

type Weather_weather struct {
    Main string         `json:"main"`
    Description string  `json:"description"`
}

type Weather_main struct {
    Temp float32        `json:"temp"`
    Humidity float32    `json:"humidity"`
    Pressue float32	`json:"pressure"`
}

type Weather_wind struct {
    Speed float32       `json:"speed"`
}

var openweather_reg = regexp.MustCompile(`(?i)\A\.wet(?:\s+|\z)(\S+)*(?:\s+|\z)([a-z]*)`) 

//posts weather
func Openweather(stored string, conn *irc.Connection) {
    location := openweather_reg.FindStringSubmatch(stored)

    if len(location) > 0{

//        var client = &http.Client{Timeout: 10 * time.Second}
        var api_url = `https://api.openweathermap.org/data/2.5/weather?units=imperial&zip=`

        fmt.Println(location)

        if location[1] != "" {
            api_url += location[1]
        } else {return}		//check map for registered zip codes
        if location[2] != "" {
            api_url += "," + location[2]
        }

        api_url += `&appid=` + OW_apikey

        fmt.Println(api_url)
/*
        resp, err := client.Get(api_url)
        Err_check(err)
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            body, err := ioutil.ReadAll(resp.Body)
            Err_check(err)
            
            var dat map[string]json.RawMessage
            err := json.Unmarshal(Test_weather, &dat)
            Err_check(err)

            var name string
            err = json.Unmarshal(dat["name"], &name)
            Err_check(err)

            var ww []Weather_weather
            err = json.Unmarshal(dat["weather"], &ww)
            Err_check(err)

            var wm Weather_main
            err = json.Unmarshal(dat["main"], &wm)
            Err_check(err)

            var wwi Weather_wind
            err = json.Unmarshal(dat["wind"], &wwi)
            Err_check(err)
        }
*/
    }

}

var youtube_linkreg = regexp.MustCompile(`(https:\/\/|https:\/\/www.)(youtube.com\/watch\?v=|youtu.be\/)(\S+)`)

//posts title and description of youtube videos
func Youtube(stored string, conn *irc.Connection) {
    url := youtube_linkreg.FindStringSubmatch(stored)

    if len(url) == 4 {
        id := url[3]

        var client = &http.Client{Timeout: 10 * time.Second}

        api_url := "https://www.googleapis.com/youtube/v3/videos?part=snippet&fields=pageInfo/totalResults,items/snippet/description,items/snippet/title&id=" +
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

            pageInfo := dat["pageInfo"].(map[string]interface{})
            totalResults := pageInfo["totalResults"].(float64)
            if totalResults != 1 {return}

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

        go Youtube(stored, conn)
        go Openweather(stored, conn)
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

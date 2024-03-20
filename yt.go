package main
//youtube feature

import (
    "regexp"
    "net/http"
    "time"
    "io/ioutil"
    "unicode/utf8"
    "encoding/json"

    "github.com/thoj/go-ircevent"
)


var youtube_linkreg = regexp.MustCompile(`(?:.*)(https:\/\/|https:\/\/www.)(youtube.com\/watch\?v=|youtu.be\/)([^\?\s]+)(?:.*)`)

//posts title and description of youtube videos
func Youtube(stored, ch string, conn *irc.Connection) {
    url := youtube_linkreg.FindStringSubmatch(stored)

    if len(url) == 4 {
        id := url[3]

        var client = &http.Client{Timeout: 10 * time.Second}

        api_url := "https://www.googleapis.com/youtube/v3/videos?part=snippet&fields=pageInfo/totalResults,items/snippet/description,items/snippet/title&id=" +
                    id + "&key=" + YT_apikey
        resp, err := client.Get(api_url)
        if err != nil {
            return
        }
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

            trundesc := []rune(desc)
            var fdesc string
            dlen := utf8.RuneCountInString(desc)

            if dlen > 80 {
                dlen = 80
                fdesc = string(trundesc[:dlen])
                fdesc += "..."
            } else {
                fdesc = string(trundesc)
            }

            title = Vowel_replace(title)
            fdesc = Vowel_replace(fdesc)
            conn.Privmsg(ch, title + " ★ " + fdesc)
        }
    }
}
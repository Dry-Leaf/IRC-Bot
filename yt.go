package main
//youtube feature

import (
    "regexp"
    "net/http"
    "time"
    "html"
    "io/ioutil"
    "unicode/utf8"
    "encoding/json"

    "github.com/thoj/go-ircevent"
)


var youtube_linkreg = regexp.MustCompile(`(?:.*)(https:\/\/|https:\/\/www.)(youtube.com\/watch\?v=|youtu.be\/)([^\?\s]+)(?:.*)`)
var youtube_searchreg = regexp.MustCompile(`(?i)\A\.yt\s(.+)`)

func result_output(api_url, ch string, conn *irc.Connection) []string {
        var client = &http.Client{Timeout: 10 * time.Second}
        resp, err := client.Get(api_url)
        if err != nil {
            return nil
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
            if totalResults < 1 {return nil}

            items := dat["items"].([]interface{})[0]
            id_sec := items.(map[string]interface{})["id"]
            var videoId string
            if id_sec != nil {videoId = id_sec.(map[string]interface{})["videoId"].(string)}

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

            fdesc = Vowel_replace(fdesc)
            return []string{title, fdesc, videoId}
        }
    return nil
}

//posts title and description of youtube videos
func Youtube(stored, ch string, conn *irc.Connection) {
    url := youtube_linkreg.FindStringSubmatch(stored)

    if len(url) == 4 {
        id := url[3]

        api_url := "https://www.googleapis.com/youtube/v3/videos?part=snippet&fields=pageInfo/totalResults,items/snippet/description,items/snippet/title&id=" +
                    id + "&key=" + YT_apikey
        ro := result_output(api_url, ch, conn)

        if ro != nil {
            conn.Privmsg(ch, Vowel_replace(html.UnescapeString(ro[0])) + " ★ " + ro[1])
        }
    }
}

func YoutubeSearch(stored, ch string, conn *irc.Connection) {
    query := youtube_searchreg.FindStringSubmatch(stored)

    if len(query) == 2 {
        query_string := query[1]

        api_url := "https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=1&fields=items/id/videoId,pageInfo/totalResults,items/snippet/description,items/snippet/title&q=" +
                    query_string + "&key=" + YT_apikey
        ro := result_output(api_url, ch, conn)
        
        if ro != nil {
            conn.Privmsg(ch, Vowel_replace(html.UnescapeString(ro[0])) + " ★ youtu.be/" + ro[2])
        }
    }
}
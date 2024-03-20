package main
//open weather feature

import (
    "regexp"
    "net/http"
    "time"
    "encoding/json"
    "io/ioutil"
    "strconv"

    "github.com/thoj/go-ircevent"
)

var weather_map = map[int]string{3: "☂", 5: "☂", 2: "☈",
    6: "⛄︎", 7: "♨", 800: "☀", 8: "☁"} 

type Weather_weather struct {
    Id int              `json:"id"`
    Main string         `json:"main"`
    Description string  `json:"description"`
}

type Weather_main struct {
    Temp float64        `json:"temp"`
    Humidity float64    `json:"humidity"`
    Pressure float64	`json:"pressure"`
}

type Weather_wind struct {
    Speed float64       `json:"speed"`
}

type Weather_sys struct {
    Country string      `json:"country"`
}

var openweather_reg = regexp.MustCompile(`(?i)\A\.wet(?:\s+|\z)(\S[^,]+(?:\z|,))(?:\s+|\z)([a-z]*)`)
var openweather_cityreg = regexp.MustCompile(`(?i)\A\D+\z`)

//posts weather
func Openweather(stored, ch string, conn *irc.Connection) {
    location := openweather_reg.FindStringSubmatch(stored)

    if len(location) > 0{

        var client = &http.Client{Timeout: 10 * time.Second}
        var api_url = `https://api.openweathermap.org/data/2.5/weather?units=imperial`

        if location[1] != "" {
            //city test
            if match := openweather_cityreg.MatchString(location[1]); match {
                api_url += "&q="
            } else { api_url += "&zip=" }
            
            api_url += location[1]
        } else {return}		//check map for registered zip codes
        if location[2] != "" {
            api_url += "," + location[2]
        }

        api_url += `&appid=` + OW_apikey
        
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

            var cod int
            err = json.Unmarshal(dat["cod"], &cod)
            Err_check(err)

            if cod != 200 {return}

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

            var ws Weather_sys
            err = json.Unmarshal(dat["sys"], &ws)
            Err_check(err)

            post_id := ww[0].Id
            if post_id != 800 {
                post_id /= 100
            }

            met_temp := (wm.Temp - 32) / 1.8
            met_speed := wwi.Speed * .44704

            weather_output := ": : " + name + ", " + ws.Country + " : : " +
                weather_map[post_id] + " " + ww[0].Description + " : : " +
                "Temperature " + strconv.FormatFloat(met_temp, 'f', 2, 32) + "C - " +
                    strconv.FormatFloat(wm.Temp, 'f', 2, 32) + "F : : " +
                "Pressure " + strconv.FormatFloat(wm.Pressure, 'f', 2, 32) + "㍱ : : " +
                "Humidity " + strconv.FormatFloat(wm.Humidity, 'f', 2, 32) + "% : : " +
                "Wind " + strconv.FormatFloat(met_speed, 'f', 2, 32) + "m/s - " + 
                    strconv.FormatFloat(wwi.Speed, 'f', 2, 32) + "mph : : " +
                "https://openweathermap.org : :"

           weather_output = Vowel_replace(weather_output)
           conn.Privmsg(ch, weather_output)
        }
    }

}
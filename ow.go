package main
//open weather feature

import (
    "regexp"
    "net/http"
    "time"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "database/sql"
    "math"
    //"fmt"

    "github.com/thoj/go-ircevent"
    _ "github.com/mattn/go-sqlite3"
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
    FeelsLike float64   `json:"feels_like"`
    Humidity int64      `json:"humidity"`
    Pressure int64	`json:"pressure"`
}

type Weather_wind struct {
    Speed float64       `json:"speed"`
}

type Weather_sys struct {
    Country string      `json:"country"`
}

var openweather_reg = regexp.MustCompile(`(?i)\A\.wet(?:\s+|\z)(\S[^,]+(?:\z|,))(?:\s+|\z)([a-z]*)`)
var openweather_cityreg = regexp.MustCompile(`(?i)\A\D+\z`)

var register_reg = regexp.MustCompile(`(?i)\A\.wet_register(?:\s+)((\S[^,]+(?:\z|,))(?:\s+|\z)([a-z]*))`)
var database_reg = regexp.MustCompile(`(?i)\A\.wet(?:\s*)(?:\z)`)
var location_reg = regexp.MustCompile(`(?i)(\S[^,]+(?:\z|,))(?:\s+|\z)([a-z]*)`)

func WetRegister(sender, stored, ch string, conn *irc.Connection) {
    location := register_reg.FindStringSubmatch(stored)

    if len(location) > 0 {
        dbconn, err := sql.Open("sqlite3", "message.db")
        Err_check(err)
        defer dbconn.Close()

        statement, err := dbconn.Prepare(`INSERT OR REPLACE INTO wet (User, Location) 
            VALUES (?, ?)`)
        Err_check(err)

        statement.Exec(sender, location[1])
        conn.Privmsg(ch, string('\u0003') + "7Registration Complete");
}}

//posts weather
func Openweather(sender, stored, ch string, conn *irc.Connection) {
    var location []string
    
    if match := database_reg.MatchString(stored); match {
        dbconn, err := sql.Open("sqlite3", "message.db")
        Err_check(err)
        defer dbconn.Close()

        statement, err := dbconn.Prepare(`SELECT Location FROM wet WHERE User = ?`)
        Err_check(err)

        var temp_location string
        err = statement.QueryRow(sender).Scan(&temp_location)
        if err == sql.ErrNoRows {
            conn.Privmsg(ch, "No location associated with: " + sender + ". Please use .wet_register to register one.");
        } else {Err_check(err)}
        
        location = location_reg.FindStringSubmatch(temp_location)
    } else {location = openweather_reg.FindStringSubmatch(stored)}

    if len(location) > 0{

        var client = &http.Client{Timeout: 10 * time.Second}
        var api_url = `https://api.openweathermap.org/data/2.5/weather?units=imperial`

        if location[1] != "" {
            //city test
            if match := openweather_cityreg.MatchString(location[1]); match {
                api_url += "&q="
            } else { api_url += "&zip=" }
            
            api_url += location[1]
        } else {return}		
        if location[2] != "" {
            api_url += location[2]
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

            var city_id int64
            err = json.Unmarshal(dat["id"], &city_id)
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

            var ws Weather_sys
            err = json.Unmarshal(dat["sys"], &ws)
            Err_check(err)

            post_id := ww[0].Id
            if post_id != 800 {
                post_id /= 100
            }

            met_temp := (wm.Temp - 32) / 1.8
            met_fl := (wm.FeelsLike - 32) / 1.8
            met_speed := wwi.Speed * .44704

            seperator := " " + string('\u0003') + "7●" + string('\u0003') + " "
            slash := " " + string('\u0003') + "7FL" + string('\u0003') + " "

            weather_output := string('\u0002') + name + ", " + ws.Country + string('\u0002') + seperator +
                strconv.FormatFloat(math.Round(met_temp), 'f', 0, 32) + "℃ - " + strconv.FormatFloat(math.Round(wm.Temp), 'f', 0, 32) + "℉" 
                
            if (math.Round(wm.FeelsLike) != math.Round(wm.Temp)) {
                weather_output += slash +
                strconv.FormatFloat(math.Round(met_fl), 'f', 0, 32) + "℃ - " + strconv.FormatFloat(math.Round(wm.FeelsLike), 'f', 0, 32) + "℉"}

            weather_output += seperator +
                weather_map[post_id] + " " + ww[0].Description + seperator +
                strconv.FormatFloat(math.Round(met_speed), 'f', 0, 32) + "㎧ - " + strconv.FormatFloat(math.Round(wwi.Speed), 'f', 0, 32) + "mph" + seperator +
                strconv.FormatInt(wm.Humidity, 10) + "%" + seperator +
                strconv.FormatInt(wm.Pressure, 10) + "㍱" + seperator

           weather_output = Vowel_replace(weather_output) + "https://openweathermap.org"
           if city_id != 0 {weather_output += "/city/" + strconv.FormatInt(city_id, 10)}

           conn.Privmsg(ch, weather_output)
        } else {conn.Privmsg(ch, string('\u0003') + "4Invalid Location");}
    }

}
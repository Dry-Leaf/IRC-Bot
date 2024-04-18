package main

import (
    "time"
    "regexp"
    "io"
    "math"
    "net/http"
    "encoding/json"
//"fmt"
    "github.com/thoj/go-ircevent"
    "github.com/IvanMenshykov/MoonPhase"
    "github.com/Lofanmi/chinese-calendar-golang/lunar"
    "github.com/hashicorp/go-set/v2"
)

var UTC_loc *time.Location

var calendar_reg = regexp.MustCompile(`(?i)\A\.cal(?:\s+|\z)`)

var phase_output = map[int]string {
    0 : "🌑︎ New Moon",
    1 : "🌒︎ Waxing Crescent",
    2 : "🌓︎ First Quarter",
    3 : "🌔︎ Waxing Gibbous",
    4 : "🌕︎ Full Moon",
    5 : "🌖︎ Waning Gibbous",
    6 : "🌗︎ Third Quarter",
    7 : "🌘︎ Waning Crescent",
    8 : "🌑︎ New Moon",
}

var zodiac_output = map[string]string {
    "aries": "♈︎",
    "taurus": "♉︎",
    "gemini": "♊︎",
    "cancer": "♋︎",
    "leo": "♌︎",
    "virgo": "♍︎",
    "libra": "♎︎",
    "scorpio": "♏︎",
    "sagittarius": "♐︎",
    "capricorn": "♑︎",
    "aquarius": "♓︎",
    "pisces": "♓︎",
}

//USA, Canada, Japan, South Korea, China, Hong Kong, France, Canda, Brazil, Austria, Australia, Belarus, Switzerland, 
//Spain, Finland, Norway, Faroe Islands, Lithuania, Czechia, Slovakia, Croatia
var country_codes = set.From[string]([]string{"US", "CA", "JP", "KR", "CN", "HK", "FR", "BR", "AT", 
    "AU", "BY", "CH", "ES", "FI", "NO", "FO", "LT", "CZ", "SK", "HR"})

type holiday_struct struct {
    Date string            `json:"date"`
    Name string            `json:"name"`
    CountryCode string     `json:"countryCode"`
}

func get_holidays(now time.Time) string {
    var client = &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get("https://date.nager.at/api/v3/NextPublicHolidaysWorldwide")
    if err != nil {
        return ""
    }
    defer resp.Body.Close()

    var holiday_output string
    if resp.StatusCode == http.StatusOK {
        body, err := io.ReadAll(resp.Body)
        Err_check(err)

        var received []holiday_struct
        err = json.Unmarshal(body, &received)
        Err_check(err)

        now_date := now.Format(time.DateOnly)
        holidays_today := make(map[string]struct{})

        for _, holiday := range received {
            if holiday.Date != now_date {continue}
            if incl := country_codes.Contains(holiday.CountryCode); !incl {continue}

            holidays_today[holiday.Name] = struct{}{}
        }

        num := len(holidays_today)
        counter := 0

        for holiday, _ := range holidays_today{
            holiday_output += holiday
            if (counter + 1) < num {holiday_output += " ★ "}
            counter += 1
    }}
    return holiday_output
}

func calendar_output(now time.Time, ch string, conn *irc.Connection) {
    //year, month, day := now.Date()
    //midnight := time.Date(year, month, day, 0, 0, 0, 0, UTC_loc).AddDate(0, 0, 1)
    cal_output := string('\u0002') + now.Format("Mon Jan 2 ") + string('\u0002')

    moon := MoonPhase.New(now)
    phase := phase_output[int(math.Floor((moon.Phase() + 0.0625 ) * 8))]
    zodiac := zodiac_output[moon.ZodiacSign()]

    lunar_calc := lunar.NewLunar(&now)
    animal := lunar_calc.Animal().Alias()

    holidays := get_holidays(now)
    if len(holidays) > 1 {holidays = "🗓 " + holidays}

    cal_output += " ★ " + phase + " ★ " + zodiac + " " + animal

    conn.Privmsg(ch, cal_output)
    time.Sleep(time.Second)
    conn.Privmsg(ch, holidays)
}

func Calendar_auto(conn *irc.Connection) {
    aux := func(now time.Time, conn *irc.Connection) {
        for _, ch := range Channels {calendar_output(now, ch, conn)}
    }

    time.Sleep(5 * time.Second)

    var err error
    UTC_loc, err = time.LoadLocation("UTC")
    Err_check(err)

    now := time.Now().UTC()
    year, month, day := now.Date()
    midnight := time.Date(year, month, day, 0, 0, 0, 0, UTC_loc).AddDate(0, 0, 1)

    aux(now, conn)

    for range time.Tick(1 * time.Minute) {
        now = time.Now().UTC()
        if now.After(midnight) {
            aux(now, conn)
            year, month, day = now.Date()
            midnight = time.Date(year, month, day, 0, 0, 0, 0, UTC_loc).AddDate(0, 0, 1)
    }}
}

func Calendar(stored, ch string, conn *irc.Connection) {
    match := calendar_reg.MatchString(stored)
    if !match {return}

    now := time.Now().UTC()
    calendar_output(now, ch, conn)
}

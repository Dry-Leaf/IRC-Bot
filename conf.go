package main

import (
    "log"

    ini "gopkg.in/ini.v1"
)

var Server string
var Channels []string
var Nickname, Username string

var YT_apikey string
var OW_apikey string
var Def_usernames = make(map[string]bool)

func blank_check(setting string, key string) {
    if setting == "" {
        log.Fatal(key + "'s value in chii.ini appears to be missing")
    }
}

func Load_conf() {
    cfg, err := ini.LoadSources(
        ini.LoadOptions{AllowBooleanKeys: true,}, "chii.ini")
    Err_check(err)

    Server = cfg.Section("IRC Settings").Key("server").String(); blank_check(Server, "server")
    Channels = cfg.Section("IRC Settings").Key("channel").Strings(" ")
    for i, ch := range Channels {
        blank_check(ch, "channel")
        Channels[i] = "#" + ch
    }

    Nickname = cfg.Section("IRC Settings").Key("nickname").String(); blank_check(Nickname, "nickname")
    Username = cfg.Section("IRC Settings").Key("username").String(); blank_check(Username, "username")

    YT_apikey = cfg.Section("API Keys").Key("youtube").String(); blank_check(YT_apikey, "youtube")
    OW_apikey = cfg.Section("API Keys").Key("openweather").String(); blank_check(OW_apikey, "openweather")

    defaults := cfg.Section("Default Usernames").KeyStrings()
    for _, u := range defaults {
        Def_usernames[u] = true
    }
}

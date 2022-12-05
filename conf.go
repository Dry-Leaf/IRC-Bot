package main

import (
    "log"

    ini "gopkg.in/ini.v1"
)

var Server, Channel string
var Nickname, Username string

var YT_apikey string

func blank_check(setting string, key string) {
    if setting == "" {
        log.Fatal(key + "'s value in chi.ini appears to be missing")
    }
}

func Load_conf() {
    cfg, err := ini.Load("chi.ini")
    Err_check(err)

    Server = cfg.Section("IRC Settings").Key("server").String(); blank_check(Server, "server")
    Channel = cfg.Section("IRC Settings").Key("channel").String(); blank_check(Channel, "channel")
    Channel = "#" + Channel

    Nickname = cfg.Section("IRC Settings").Key("nickname").String(); blank_check(Channel, "nickname")
    Username = cfg.Section("IRC Settings").Key("username").String(); blank_check(Channel, "username")

    YT_apikey = cfg.Section("API Keys").Key("youtube").String(); blank_check(Channel, "youtube")
}

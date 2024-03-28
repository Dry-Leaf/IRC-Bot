# Chii

An irc bot written in golang.

## Features

- Posts youtube video's title and description

## Compile Instructions

git clone https://gitgud.io/nvtelen/chii

cd chii 

go mod init modules

go mod tidy 

go build -o chii *.go

## Comands

`Look up weather`

.wet

.wet New York

.wet Paris, US

.wet 120-0015, JP

`Register Weather location`

.wet_register *location*

`Trivia`
.triv

.triv 10

.striv

`Send private message to someone the next time they enter the channel`

.tell *recipient* *message*

`Look up words in Urban Dictionary`

.ud word

`8ball`

.8ball

`Fortune Cookie`

.fortune

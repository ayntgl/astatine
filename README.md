# Astatine

[![Go Reference](https://pkg.go.dev/badge/github.com/bwmarrin/discordgo.svg)](https://pkg.go.dev/github.com/bwmarrin/discordgo) [![Go Report Card](https://goreportcard.com/badge/github.com/bwmarrin/discordgo)](https://goreportcard.com/report/github.com/bwmarrin/discordgo) [![Build Status](https://travis-ci.com/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.com/bwmarrin/discordgo) [![Discord Gophers](https://img.shields.io/badge/Discord%20Gophers-%23discordgo-blue.svg)](https://discord.gg/golang) [![Discord API](https://img.shields.io/badge/Discord%20API-%23go_discordgo-blue.svg)](https://discord.com/invite/discord-api)

A powerful, versatile, and efficient Discord API library.

## Getting Started

### Installation

```
go get github.com/ayntgl/astatine
```

### Usage

```go
package main

import (
    "os"
    "os/signal"
    "fmt"

    "github.com/ayntgl/astatine"
)

func main() {
    token := os.Getenv("DISCORD_TOKEN")
    session := astatine.New(token)

    err := session.Open()
    if err != nil {
        panic(err)
    }

    fmt.Println("Press Ctrl+C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, os.Interrupt)
    <-sc
}
```

## List of Discord APIs

See [this chart](https://abal.moe/Discord/Libraries.html) for a feature 
comparison and list of other Discord API libraries.

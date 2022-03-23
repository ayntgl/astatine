# Astatine

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

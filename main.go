package main

import (
    "log"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    log.Print("Port :" + port)

    a := App{}
    a.Initialize("root", "", "movies-api")
    a.Run(":" + port)
}

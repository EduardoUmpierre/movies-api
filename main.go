package main

func main() {
    addr, err := determineListenAddress()
    if err != nil {
        log.Fatal(err)
    }

    a := App{}
    a.Initialize("root", "", "movies-api")
    a.Run(addr)
}

func determineListenAddress() (string, error) {
    port := os.Getenv("PORT")
    if port == "" {
        return "", fmt.Errorf("$PORT not set")
    }
    return ":" + port, nil
}

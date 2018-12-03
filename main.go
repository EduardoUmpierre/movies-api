package main

func main() {
    a := App{}
    a.Initialize("root", "", "movies-api")
    a.Run(":8080")
}

package main

import (
  "fmt"
  "github.com/shavit/gochat/server"
  "github.com/shavit/gochat/cmd"
)

func startServer()  {
  var s server.Server = server.NewServer()
  s.Start()
}

func startClient(){
  cmd.RunClient()
}

func startServerOrClient(){
  var input string

  println(`
    Are you a server or a client?
    1 - Server
    2 - Client
  `)
  fmt.Scanln(&input)

  switch input {
  case "1":
    startServer()
  case "2":
    startClient()
  default:
    startServer()
  }
}

func main()  {
  startServerOrClient()
}

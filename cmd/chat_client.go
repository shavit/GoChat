package cmd

import (
  "os"
  "os/signal"
  "net"
  "io"
  "syscall"
  "fmt"
  "bufio"
)

type ChatClient interface {
  Dial(network, address string) (err error)
  Close()
  setUsername()
  handleInput()
  echoServer(done chan struct{})
}

type chatClient struct {
  conn net.Conn
}

func NewChatClient() ChatClient {
  return new(chatClient)
}

func (client *chatClient) Dial(network, address string) (err error){
  client.conn, err = net.Dial(network, address)
  return err
}

func (client *chatClient) Close(){
  client.conn.Close()
}

func (client *chatClient) setUsername() (){
  var username string

  for len(username) == 0 {
    fmt.Println("Choose a username")
    fmt.Scan(&username)
  }

  io.WriteString(client.conn, "USERNAME:"+username)
}

func (client *chatClient) handleInput(){
  var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
  for scanner.Scan(){
    io.WriteString(client.conn, "MESSAGE:"+scanner.Text())
  }
}

func (client *chatClient) echoServer(done chan struct{}){
  for {
    buf := make([]byte, 12000)
    n, err := client.conn.Read(buf)
    if err != nil {
      print("---> Connection closed by the server. ")
      break
    }

    fmt.Println(string(buf[0:n]))
    buf = nil
  }

  println("Disconnecting")
  done <- struct{}{}
}


func RunClient(){
  var err error
  var ch chan os.Signal = make(chan os.Signal)
  var client ChatClient = NewChatClient()
  var done chan struct{} = make(chan struct{})

  if err = client.Dial("tcp", "gochat_server:2400"); err != nil {
    panic(err)
  }
  defer client.Close()

  go func() {
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    <-ch
    println("\n\nDisconnecting from the server\n")
    os.Exit(0)
  }()

  go client.handleInput()
  go client.echoServer(done)
  client.setUsername()
  <-done
}

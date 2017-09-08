package server

import (
  "os"
  "os/signal"
  "io"
  "errors"
  "syscall"
  "net"
  "github.com/shavit/gochat/chatroom"
)

type Server interface {
  Start() (err error)
  routeClient(conn net.Conn)
  printRooms(conn io.ReadWriteCloser)
  AddRoom(name string)
  RemoveRoom(name string)
  GetRoomsKeys() (keys []string)
}

type server struct {
  rooms map[string]chatroom.ChatRoom
}

func NewServer() (Server){
  return &server{
    rooms: make(map[string]chatroom.ChatRoom),
  }
}

func (s *server) Start() (err error){
  var hostname string = os.Getenv("SERVER_HOSTNAME")
  var ln net.Listener
  var conn net.Conn
  var ch chan os.Signal = make(chan os.Signal)

  go func() {
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    // Block and wait for a signal to exit the program
    <-ch
    println("\n\nSending exit notificaiton to all clients\n")
    ln.Close()
    os.Exit(0)
  }()


  ln, err = net.Listen("tcp", hostname+":2400")
  if err != nil {
    return errors.New("Error listening on port 2400")
  }
  defer ln.Close()
  println("---> Listening on tcp://"+hostname+":2400")

  for {
    conn, err = ln.Accept()
    if err != nil {
      println(err)
      errors.New("Error accepting messages")
    }

    go s.routeClient(conn)
  }

  return err
}


func (s *server) routeClient(conn net.Conn){

  println("Available rooms:")
  s.printRooms(conn)

  // TODO: Choose the room based on the user input
  s.rooms["lobby"].AddUser(conn)
}

func (s *server) printRooms(conn io.ReadWriteCloser){
  if len(s.rooms) < 1 {
    s.AddRoom("lobby")
  }

  for i, _ := range(s.rooms){
    println(i, ":", s.rooms[i].String())
  }
}

func (s *server) AddRoom(key string){
  s.rooms[key] = chatroom.NewRoom(key)
}

func (s *server) RemoveRoom(key string){
  delete(s.rooms, key)
}

func (s *server) GetRoomsKeys() (keys []string){
  for k := range s.rooms {
    keys = append(keys, k)
  }

  return keys
}

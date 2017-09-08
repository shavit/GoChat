package chatroom

import (
  "io"
  "sync"
  "time"
  "errors"
)

type ChatRoom interface {
  String() (description string)
  Close()
  listen(u *user)
  parse(msg []byte) (message *incomingMessage, err error)
  Dispatch(u *user, message *incomingMessage) (err error)
  broadcast(msg []byte, from *user)
  AddUser(u io.ReadWriteCloser)
  Users() (users []user)
  removeUser(wc chan<- string)
}

type user struct {
  Id int
  Name string
  RoomName string
  rwc io.ReadWriteCloser
}

type incomingMessage struct {
  option string
  body []byte
}

type chatroom struct {
  Name string
  users map[chan<- string]*user
  sync.RWMutex
}

func NewRoom(name string) ChatRoom {
  return &chatroom{
    Name: name,
    users: make(map[chan<- string]*user),
  }
}

func (r *chatroom) String() (description string){
  return r.Name
}

func (r *chatroom) Close(){
  println("---> Closing room", r.Name)
}

func (r *chatroom) listen(u *user){
  var err error
  var conn io.ReadWriteCloser = u.rwc
  var buf []byte
  var n int
  var message *incomingMessage

  defer conn.Close()
  L:
  for {
    buf = make([]byte, 12000)
    n, err = conn.Read(buf)
    switch err {
    case io.EOF:
      println(time.Now().String(), u.Name, "has left")
      break L
    case nil:
      break
    default:
      panic(err)
    }

    message, err = r.parse(buf[0:n])
    if err != nil {
      panic(err)
    }

    err = r.Dispatch(u, message)
    if err != nil {
      panic(err)
    }

  }
}

func (r *chatroom) parse(msg []byte) (message *incomingMessage, err error){
  var n int = len(msg)

  for i := range msg {
    if string(msg[i:i+1]) == ":"{
      return &incomingMessage{
        option: string(msg[0:i]),
        body: msg[i+1:n],
      }, err
    }
  }

  println(string(msg))

  return message, errors.New("Invalid message format")
}

func (r *chatroom) Dispatch(u *user, message *incomingMessage) (err error){
  var conn io.ReadWriteCloser = u.rwc

  switch message.option {
  case "USERNAME":
    u.Id = len(r.Users())
    u.Name = string(message.body)
    io.WriteString(conn, "\n\n          🙋‍♂️ Welcome to the chat server 🙋\n")
    io.WriteString(conn, "          🙋‍♂️ You are connected as: "+string(message.body)+" 🙋\n\n")
    println(time.Now().String(), u.Name, "has joined")
    break

  case "MESSAGE":
    r.broadcast(message.body, u)
    break

  default:
    err = errors.New("Unhandled message option")
  }

  return err
}

func (r *chatroom) broadcast(msg []byte, from *user){
  var u *user

  print("---> Broadcasting to ", len(r.users), " users. ")
  print(from.Name+": ")
  println(string(msg))

  r.RLock()
  defer r.RUnlock()
  for _, u = range r.users {
    go func(user *user) {
      io.WriteString(user.rwc, from.Name+": "+string(msg))
    }(u)
  }
}

func (r *chatroom) AddUser(rwc io.ReadWriteCloser){
  var writeChannel chan<- string = make(chan string)

  r.Lock()
  r.users[writeChannel] = &user{
    Id: len(r.users),
    rwc: rwc,
  }
  r.Unlock()

  r.listen(r.users[writeChannel])
  defer r.removeUser(writeChannel)
}

func (r *chatroom) Users() (users []user){
  return users
}

func (r *chatroom) removeUser(wc chan<- string){
  println("---> Removing user", wc)
  r.Lock()
  close(wc)
  delete(r.users, wc)
  r.Unlock()
}

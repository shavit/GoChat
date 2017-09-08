package chatroom

import (
  "testing"
)

type mConn struct {}
func (c mConn) Read([]byte) (int, error){
  return 0, nil
}
func (c mConn) Write([]byte) (int, error){
  return 0, nil
}
func (c mConn) Close() (error){
  return nil
}

func TestCreateRoom(t *testing.T){
  var room ChatRoom = NewRoom("Test Room")
  var chatRoom *chatroom
  var ok bool

  if room == nil {
    t.Error("Could not create a room")
  }

  chatRoom, ok = room.(*chatroom)
  if ok != true {
    t.Error("Error turning ChatRoom into chatroom")
  }

  if chatRoom.Name != "Test Room" {
    t.Error("Error setting a name for a new chatroom")
  }
}

func TestParseMessage(t *testing.T){
  var room ChatRoom = NewRoom("Test Room")
  var err error
  var msg []byte
  var m *incomingMessage

  msg = []byte("New Message")
  m, err = room.parse(msg)
  if m != nil {
    t.Error("Should not parse invalid messages.")
  }

  if err.Error() != "Invalid message format"{
    println(err)
    t.Error("Should catch errors")
  }

  msg = []byte("MESSAGE:A message from the user")
  m, err = room.parse(msg)
  if m.option != "MESSAGE"{
    t.Error("Should parse an option from an incoming message")
  }
  if string(m.body) != "A message from the user"{
    t.Error("Should return the message body")
  }

  msg = []byte("USERNAME:Mike")
  m, err = room.parse(msg)
  if m.option != "USERNAME"{
    t.Error("Should parse an option from an incoming message")
  }
  if string(m.body) != "Mike"{
    t.Error("Should return the message body")
  }
}

package server

import (
  "testing"
  "io"
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

func TestServerCreate(t *testing.T){
  var server Server = NewServer()

  if server == nil {
    t.Error("Error creating a server")
  }
}

func TestCreateDefaultRoom(t *testing.T){
  var server Server = NewServer()
  var count int = len(server.GetRoomsKeys())
  var name string
  var conn io.ReadWriteCloser = mConn{}

  if count > 0 {
    t.Error("Server should not have rooms, got", count)
  }

  server.printRooms(conn)
  count = len(server.GetRoomsKeys())
  name = server.GetRoomsKeys()[0]
  if count != 1 || name != "lobby" {
    t.Error("Server should have a default lobby room, got", name)
  }
}

func TestAddRoom(t *testing.T){
  var server Server = NewServer()
  var names []string

  server.AddRoom("Room 1")
  server.AddRoom("Room 2")

  if len(server.GetRoomsKeys()) != 2 {
    t.Error("No rooms")
  }

  names = server.GetRoomsKeys()
  if names[0] != "Room 1" || names[1] != "Room 2"{
    t.Error("Error adding a room to server")
  }
}

func TestRemoveRoom(t *testing.T){
  var server Server = NewServer()
  var count int
  var name string

  server.AddRoom("Lobby")
  server.AddRoom("Random")
  server.RemoveRoom("Lobby")

  count = len(server.GetRoomsKeys())
  name = server.GetRoomsKeys()[0]
  if count != 1 || name != "Random" {
    t.Error("Server should have a default lobby room, got", name)
  }
}

package core

import (
    "fmt"
    "github.com/gorilla/websocket"
)

type Room struct {
    ID         string
    Sessions   map[string]*StreamSession
    CallerConn *websocket.Conn
}

type StreamSession struct {
    ID                  string
    Offer               string
    Answer              string
    CallerIceCandidates []string
    CalleeIceCandidates []string
    CallerConn          *websocket.Conn
    CalleeConn          *websocket.Conn
}

var roomMap = make(map[string]*Room)

func GetRoom(id string) *Room {
    return roomMap[id]
}

func NewRoom(callerConn *websocket.Conn) *Room {
    room := Room{
        ID:         newRoomID(),
        Sessions:   make(map[string]*StreamSession),
        CallerConn: callerConn,
    }
    roomMap[room.ID] = &room
    return &room
}

func newRoomID() string {
    id := GetRandomName(0)
    for GetRoom(id) != nil {
        id = GetRandomName(0)
    }
    return id
}

func RemoveRoom(id string) {
    roomMap[id] = nil
}

func (room *Room) GetSession(id string) *StreamSession {
    return room.Sessions[id]
}

func (room *Room) NewSession(calleeConn *websocket.Conn) *StreamSession {
    session := StreamSession{
        ID:                  room.newSessionID(),
        CallerIceCandidates: []string{},
        CalleeIceCandidates: []string{},
        CallerConn:          room.CallerConn,
        CalleeConn:          calleeConn,
    }
    room.Sessions[session.ID] = &session
    return &session
}

func (room *Room) newSessionID() string {
    id := fmt.Sprintf("%s$%s", room.ID, GetRandomName(0))
    for GetRoom(id) != nil {
        id = fmt.Sprintf("%s$%s", room.ID, GetRandomName(0))
    }
    return id
}

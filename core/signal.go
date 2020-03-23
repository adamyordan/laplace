package core

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

type WSMessage struct {
    SessionID string
    Type      string
    Value     string
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func GetHttp() *http.ServeMux {
    server := http.NewServeMux()

    server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("files/static"))))

    server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "files/main.html")
    })

    server.HandleFunc("/ws_serve", func(writer http.ResponseWriter, request *http.Request) {
        conn, _ := upgrader.Upgrade(writer, request, nil)
        room := NewRoom(conn)
        _ = conn.WriteJSON(WSMessage{
            SessionID: "",
            Type:      "newRoom",
            Value:     room.ID,
        })

        go func(r *Room) {
            //noinspection ALL
            defer room.CallerConn.Close()
            for {
                var msg WSMessage
                if err := conn.ReadJSON(&msg); err != nil {
                    log.Println("websocketError.", err)
                    return
                }
                log.Println(msg)
                s := room.GetSession(msg.SessionID)
                if s == nil {
                    log.Println("session nil.", msg.SessionID)
                }
                if msg.Type == "addCallerIceCandidate" {
                    s.CallerIceCandidates = append(s.CallerIceCandidates, msg.Value)
                } else if msg.Type == "gotOffer" {
                    s.Offer = msg.Value
                }
                _ = s.CalleeConn.WriteJSON(msg)
            }
        }(room)
    })

    server.HandleFunc("/ws_connect", func(writer http.ResponseWriter, request *http.Request) {
        conn, _ := upgrader.Upgrade(writer, request, nil)

        ids, ok := request.URL.Query()["id"]
        if !ok || ids[0] == "" {
            return
        }

        room := GetRoom(ids[0])
        session := room.NewSession(conn)

        _ = room.CallerConn.WriteJSON(WSMessage{
            SessionID: session.ID,
            Type:      "newSession",
            Value:     session.ID,
        })
        _ = conn.WriteJSON(WSMessage{
            SessionID: session.ID,
            Type:      "newSession",
            Value:     session.ID,
        })

        go func(s *StreamSession) {
            //noinspection ALL
            defer s.CalleeConn.Close()
            for {
                var msg WSMessage
                if err := conn.ReadJSON(&msg); err != nil {
                    log.Println("websocketError.", err)
                    return
                }
                log.Println(msg)
                if msg.SessionID == s.ID {
                    if msg.Type == "addCalleeIceCandidate" {
                        s.CalleeIceCandidates = append(s.CalleeIceCandidates, msg.Value)
                    } else if msg.Type == "gotAnswer" {
                        s.Answer = msg.Value
                    }
                    _ = s.CallerConn.WriteJSON(msg)
                }
            }
        }(session)
    })

    return server
}

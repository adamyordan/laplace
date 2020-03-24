package core

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "time"
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

func sendHeartBeatWS(ticker *time.Ticker, conn *websocket.Conn, quit chan struct{}) {
    for {
        select {
        case <- ticker.C:
            _ = conn.WriteJSON(WSMessage{
                Type: "beat",
            })
        case <- quit:
            log.Println("heartbeat stopped")
            return
        }
    }
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
        if err := conn.WriteJSON(WSMessage{
            SessionID: "",
            Type:      "newRoom",
            Value:     room.ID,
        }); err != nil {
            log.Println("newSessionWriteJsonError.", err)
            return
        }

        go func(r *Room) {
            ticker := time.NewTicker(10 * time.Second)
            quit := make(chan struct{})
            defer func() {
                ticker.Stop()
                _ = room.CallerConn.Close()
                close(quit)
                RemoveRoom(r.ID)
                for sID, s := range r.Sessions {
                    _ = s.CalleeConn.WriteJSON(WSMessage{
                        Type: "roomClosed",
                        SessionID: sID,
                    })
                }
            }()

            go sendHeartBeatWS(ticker, conn, quit)

            //noinspection ALL
            defer room.CallerConn.Close()
            for {
                var msg WSMessage
                if err := room.CallerConn.ReadJSON(&msg); err != nil {
                    log.Println("websocketError.", err)
                    return
                }
                //log.Println(msg)
                s := room.GetSession(msg.SessionID)
                if s == nil {
                    log.Println("session nil.", msg.SessionID)
                }
                if msg.Type == "addCallerIceCandidate" {
                    s.CallerIceCandidates = append(s.CallerIceCandidates, msg.Value)
                } else if msg.Type == "gotOffer" {
                    s.Offer = msg.Value
                }
                if err := s.CalleeConn.WriteJSON(msg); err != nil {
                    log.Println("serveEchoWriteJsonError.", err)
                }
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
        if room == nil {
            _ = conn.WriteJSON(WSMessage{
                Type: "roomNotFound",
            })
            return
        }
        session := room.NewSession(conn)

        if err := room.CallerConn.WriteJSON(WSMessage{
            SessionID: session.ID,
            Type:      "newSession",
            Value:     session.ID,
        }); err != nil {
            log.Println("callerWriteJsonError.", err)
            return
        }

        if err := conn.WriteJSON(WSMessage{
            SessionID: session.ID,
            Type:      "newSession",
            Value:     session.ID,
        }); err != nil {
            log.Println("calleeWriteJsonError.", err)
            return
        }

        go func(s *StreamSession) {
            //noinspection ALL
            defer s.CalleeConn.Close()
            for {
                var msg WSMessage
                if err := conn.ReadJSON(&msg); err != nil {
                    log.Println("websocketError.", err)
                    return
                }
                //log.Println(msg)
                if msg.SessionID == s.ID {
                    if msg.Type == "addCalleeIceCandidate" {
                        s.CalleeIceCandidates = append(s.CalleeIceCandidates, msg.Value)
                    } else if msg.Type == "gotAnswer" {
                        s.Answer = msg.Value
                    }
                    if err := s.CallerConn.WriteJSON(msg); err != nil {
                        log.Println("connectEchoWriteJsonError.", err)
                    }
                }
            }
        }(session)
    })

    return server
}

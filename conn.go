// Copyright 2013 Clustertech Limited. All rights reserved.
//
// Author: jackeychen (jackeychen@clustertech.com)
package socket

import (
  "encoding/json"
  "github.com/gorilla/websocket"
)

type Conn struct {
  c        *websocket.Conn
  sessions map[string]interface{}
  id       string

  hub  *Hub
  send chan []byte
}

func (conn *Conn) ID() string {
  return conn.id
}

func (conn *Conn) Release() {
  conn.hub.Release(conn)
}

func (conn *Conn) readProcess() {
  defer conn.Release()
  for {
    _, message, err := conn.c.ReadMessage()
    if err != nil {
      break
    }
    m := make(map[string]string)
    err = json.Unmarshal(message, &m)
    if err != nil {
      break
    }

    e := &Event{Name: m["name"], Msg: m["msg"], Conn: conn}
    fn, ok := conn.hub.handlers[e.Name]
    if ok {
      fn(e)
    }
  }
}

func (conn *Conn) writeProcess() {
  defer conn.Release()
  for {
    select {
    case message, ok := <-conn.send:
      if !ok {
        conn.c.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      if err := conn.c.WriteMessage(websocket.TextMessage, message); err != nil {
        return
      }
    }
  }
}

func (conn *Conn) GetHub() *Hub {
  return conn.hub
}

func (conn *Conn) Process() {
  go conn.readProcess()
  conn.writeProcess()
}

func (conn *Conn) Send(e string, msg string) {
  m := make(map[string]string)
  m["name"] = e
  m["msg"] = msg
  data, _ := json.Marshal(m)
  conn.send <- data
}

func (conn *Conn) Broadcast(e string, msg string) {
  for _, c := range conn.hub.conns {
    c.Send(e, msg)
  }
}

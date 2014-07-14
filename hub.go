// Copyright 2013 Clustertech Limited. All rights reserved.
//
// Author: jackeychen (jackeychen@clustertech.com)
package socket

import (
  "crypto/rand"
  "github.com/gorilla/websocket"
  "net/http"
  "sync"
)

const (
  alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  idLength = 20
)

func randString(n int) string {
  var bytes = make([]byte, n)
  rand.Read(bytes)
  for i, b := range bytes {
    bytes[i] = alphanum[b%byte(len(alphanum))]
  }
  return string(bytes)
}

type Event struct {
  Name string
  Msg  string
  Conn *Conn
}

type Hub struct {
  conns    map[string]*Conn
  ch       chan *Event
  handlers map[string]func(e *Event)
  upgrader *websocket.Upgrader

  locker *sync.RWMutex

  OnRelease func(conn *Conn)
}

func NewHub(upgrader *websocket.Upgrader) *Hub {
  conns := make(map[string]*Conn)
  ch := make(chan *Event)
  handlers := make(map[string]func(e *Event))
  locker := new(sync.RWMutex)
  return &Hub{conns: conns, ch: ch, handlers: handlers, upgrader: upgrader,
    locker: locker}
}

func (hub *Hub) On(e string, f func(e *Event)) {
  hub.handlers[e] = f
}

func (hub *Hub) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
  conn, err := hub.upgrader.Upgrade(w, r, nil)
  if err != nil {
    return nil, err
  }
  id := randString(idLength)
  sessions := make(map[string]interface{})
  ch := make(chan []byte)

  hub.locker.Lock()
  defer hub.locker.Unlock()
  hub.conns[id] = &Conn{conn, sessions, id, hub, ch, new(sync.RWMutex)}
  return hub.conns[id], nil
}

func (hub *Hub) Release(conn *Conn) {
  hub.locker.Lock()
  defer hub.locker.Unlock()
  id := conn.ID()
  _, ok := hub.conns[id]
  if ok {
    if (hub.OnRelease != nil) {
      hub.OnRelease(conn)
    }
    delete(hub.conns, id)
    close(conn.send)
  }
}

func (hub *Hub) Get(id string) *Conn {
  hub.locker.RLock()
  defer hub.locker.RUnlock()
  return hub.conns[id]
}

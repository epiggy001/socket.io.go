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

type event struct {
  name string
  msg  string
  conn *Conn
}

type Hub struct {
  conns    map[string]*Conn
  ch       chan *event
  handlers map[string]func(e *event)
  upgrader *websocket.Upgrader

  locker *sync.RWMutex
}

func NewHub(upgrader *websocket.Upgrader) *Hub {
  conns := make(map[string]*Conn)
  ch := make(chan *event)
  handlers := make(map[string]func(e *event))
  locker := new(sync.RWMutex)
  return &Hub{conns, ch, handlers, upgrader, locker}
}

func (hub *Hub) On(e string, f func(e *event)) {
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
  hub.conns[id] = &Conn{conn, sessions, id, hub, ch}
  return hub.conns[id], nil
}

func (hub *Hub) Release(conn *Conn) {
  id := conn.ID()
  delete(hub.conns, id)
  close(conn.send)
}

func (hub *Hub) Get(id string) *Conn {
  hub.locker.RLock()
  defer hub.locker.RUnlock()
  return hub.conns[id]
}

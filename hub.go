// Copyright 2013 Clustertech Limited. All rights reserved.
//
// Author: jackeychen (jackeychen@clustertech.com)
package socket

import (
  "github.com/gorilla/websocket"
  "net/http"

  "crypto/rand"
)

const (
  alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  idLength = 20
)

func randString(n int) string {
  var bytes = make([]byte, n)
  rand.Read(bytes)
  for i, b := range bytes {
    bytes[i] = alphanum[b % byte(len(alphanum))]
  }
  return string(bytes)
}

type event struct {
  name string
  msg  string
  conn    *Conn
}

type Hub struct {
  conns map[string]*Conn
  ch chan *event
  handlers map[string]func(e *event)
  upgrader *websocket.Upgrader
}

func NewHub(upgrader *websocket.Upgrader) *Hub {
  conns := make(map[string]*Conn)
  ch := make(chan *event)
  handlers := make(map[string]func(e *event))
  return &Hub{conns, ch, handlers, upgrader}
}

func (hub *Hub) Upgarde(w http.ResponseWriter, r *http.Request) (*Conn, error) {
  conn, err := hub.upgrader.Upgrade(w, r, nil)
  if err != nil {
    return nil, err
  }
  id := randString(idLength)
  sessions := make(map[string]interface{})
  return &Conn{conn, sessions, id, hub}, nil
}

// Copyright 2013 Clustertech Limited. All rights reserved.
//
// Author: jackeychen (jackeychen@clustertech.com)
package socket

import (
  "github.com/gorilla/websocket"
  "net/http"
)

type event struct {
  name string
  msg  string
  c    *Conn
}

type hub struct {
  conns map[string]*Conn
  ch chan *event
  handlers map[string]func(e *event)
  upgrader *websocket.Upgrader
}

func NewHub(upgrader *websocket.Upgrader) *hub {
  conns := make(map[string]*Conn)
  ch := make(chan *event)
  handlers := make(map[string]func(e *event))
  return &hub{conns, ch, handlers, upgrader}
}

func (h *hub) Upgarde(w http.ResponseWriter, r *http.Request) *Conn {
  return nil
}

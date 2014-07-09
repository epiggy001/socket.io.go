// Copyright 2013 Clustertech Limited. All rights reserved.
//
// Author: jackeychen (jackeychen@clustertech.com)
package socket

import (
  "github.com/gorilla/websocket"
)

type Conn struct {
  *websocket.Conn
  sessions map[string]interface{}
  id string

  h *hub
}

package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
)

func NewSocketIOServer() *socketio.SocketIOServer {
	sio := socketio.NewSocketIOServer(&socketio.Config{})

	sio.Of("/diskevent").On("connect", func(ns *socketio.NameSpace) {
		go func(ns *socketio.NameSpace) {
			sub := trapTopic.Subscribe()
			defer trapTopic.Unsubscribe(sub)
			for {
				e := <-sub

				bytes, err := json.Marshal(e)
				if err != nil {
					continue
				}

				err = ns.Emit("diskevent", string(bytes))
				//err := ns.Emit("event", e)
				if err != nil {
					return
				}
			}
		}(ns)
	})

	sio.Of("/statistics").On("connect", func(ns *socketio.NameSpace) {
		go func(ns *socketio.NameSpace) {
			sub := statTopic.Subscribe()
			defer statTopic.Unsubscribe(sub)
			for {
				stat := <-sub
				/*m := make(map[string]interface{}){"status": "ok", "sample": stat}
				bytes, err := json.Marshal(m)
				if err != nil {
					continue
				}

				err = ns.Emit("statistics", bytes)*/
				//m := map[string]interface{}{"status": "ok", "sample": stat}

				//ms,err := json.Marshal(m)
				err := ns.Emit("statistics", stat)
				if err != nil {
					return
				}
			}
		}(ns)
	})

	return sio
}

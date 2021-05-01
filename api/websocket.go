package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {

	return true
}}

type client struct {
	ws *WebSocket
	outgoing chan []byte

	conn *websocket.Conn
}

func (cl *client) run() {
	go cl.reader()
	go cl.writer()
}

func (cl *client) writer() {
	for {
		msg, ok := <- cl.outgoing
		if !ok {
			log.Warn().Msg("Error reading outgoing channel")
			break
		}
		err := cl.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

// We pretty much just need this to respond to close messages and stuff
func (cl *client) reader() {
	defer func() {
		cl.ws.unregister(cl)
		cl.conn.Close()
	}()
	for {
		_, _, err := cl.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

type WebSocket struct {
	clients map[*client]bool

	addclient chan *client
	removeclient chan *client

	acnodeListener *acnode.HandlerListener
}

func CreateWebsockerHandler(handler *acnode.ACNodeHandler) *WebSocket {
	ws := WebSocket{
				clients: make(map[*client]bool),
				addclient: make(chan *client),
				removeclient: make(chan *client),
				acnodeListener: acnode.CreateHandlerChangeListener(handler, "Websocket"),
			}
	go ws.process()
	ws.acnodeListener.SetOnNodeAddedHandler(ws.onACNodeUpdate)
	ws.acnodeListener.SetOnNodeChangedHandler(ws.onACNodeUpdate)
	return &ws
}

func (ws *WebSocket) register(c *client) {
	ws.addclient <- c
}

func (ws *WebSocket) unregister(c *client) {
	ws.removeclient <- c
}

func (ws *WebSocket) onACNodeUpdate(node acnode.ACNode) {
	data,_ := json.Marshal(node.GetAPIRecord())
	go ws.Send(data)
}

func (ws *WebSocket) process() {
	for {
		select {
			case c := <- ws.addclient:
				ws.clients[c] = true
			case c := <- ws.removeclient:
				if _, ok := ws.clients[c];ok {
					delete(ws.clients, c)
					close(c.outgoing)
				}
		}
	}
}

func (ws *WebSocket) Send(msg []byte) {
	for cl := range ws.clients {
		cl.outgoing <- msg
	}
}

func (ws *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok,_ := auth.CheckAuthAPI(w, r); !ok {
		w.WriteHeader(401)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("Error upgrading websocket connection")
		w.WriteHeader(500)
		return
	}

	c := client {
		ws: ws,
		outgoing: make(chan []byte),
		conn: conn,
	}
	ws.register(&c)
	c.run()
}
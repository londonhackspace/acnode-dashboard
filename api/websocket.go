package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const (
	writeTimeout = 10 * time.Second

	// Be fairly aggressive with the pings since we don't receive anything else
	// and theoretically outoging messages might be occasional
	pingPeriod = 10 * time.Second
	pongWait = 30 * time.Second
)

var (
	clientCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_client_count",
		Help: "Number of clients connected",
	})
	connectionCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_connection_count",
		Help: "Total number of connections",
	})
	disconnectionCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_disconnection_count",
		Help: "Total number of disconnections",
	})
	messageSentCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_message_sent_count",
		Help: "Number of messages sent to clients",
	})
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
	log.Info().Str("RemoteAddr", cl.conn.RemoteAddr().String()).Msg("New Websocket Connection")
	go cl.reader()
	go cl.writer()
}

func (cl *client) writer() {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
			case msg, ok := <- cl.outgoing:
				cl.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				if !ok {
					// channel was closed
					cl.conn.WriteMessage(websocket.CloseMessage, []byte{})
					cl.ws.unregister(cl)
					cl.conn.Close()
					return
				}
				messageSentCounter.Inc()
				err := cl.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Err(err).Msg("Error writing to websocket connection")
					cl.ws.unregister(cl)
					cl.conn.Close()
					return
				}
			case <- ticker.C:
				cl.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				if err := cl.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Err(err).Msg("Error writing ping to websocket")
					return
				}
		}
	}
}

// We pretty much just need this to respond to close messages and stuff
func (cl *client) reader() {
	connectionCounter.Inc()
	defer func() {
		cl.ws.unregister(cl)
		cl.conn.Close()
		disconnectionCounter.Inc()
	}()

	cl.conn.SetReadDeadline(time.Now().Add(pongWait))
	cl.conn.SetPongHandler(func(string) error {
		cl.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

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
				clientCounter.Inc()
			case c := <- ws.removeclient:
				if _, ok := ws.clients[c];ok {
					clientCounter.Dec()
					delete(ws.clients, c)
					close(c.outgoing)
					log.Info().Msg("Websocket Connection Removed")
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
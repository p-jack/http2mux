package wsmux

import "context"
import "log"
import "net/http"
import "sync"
import "github.com/gorilla/websocket"

type Mux struct {
	auth Auth
  await *sync.WaitGroup
	conns map[string]*conns
	mutex *sync.Mutex
	pubSub PubSub
	server *http.Server
  upgrader websocket.Upgrader
}

func New(cfg Config, auth Auth, pubSub PubSub) *Mux {
  result := &Mux{}
	handler := http.NewServeMux()
  handler.HandleFunc(cfg.Endpoint, result.handle)
	result.auth = auth
	result.await = &sync.WaitGroup{}
  result.await.Add(1)
	result.conns = map[string]*conns{}
	result.mutex = &sync.Mutex{}
	result.pubSub = pubSub
	pubSub.DeliverTo(result.deliver)
	result.server = &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}
  result.upgrader = websocket.Upgrader{}
  return result
}

func (m *Mux) start() {
  defer m.await.Done()
  log.Printf("Mux Start: %s", m.server.Addr)
  err := m.server.ListenAndServe()
  log.Printf("ListenAndServe: %v", err)
}

func (m *Mux) Start() {
  go m.start()
}

func (m *Mux) Stop() {
  log.Printf("Mux Stop: %s", m.server.Addr)
  m.server.Shutdown(context.Background())
}

func (m *Mux) connsFor(topic string) *conns {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.conns[topic]
}

func (m *Mux) createConns(topic string) *conns {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := m.conns[topic]
	if result == nil {
		result = newConns(topic)
		m.conns[topic] = result
	}
	return result
}

func (m *Mux) removeConn(id int, conns *conns) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if (conns.Remove(id) == 0) {
		delete(m.conns, conns.topic)
		m.pubSub.Unsubscribe(conns.topic)
	}
}

func (m *Mux) handle(w http.ResponseWriter, r *http.Request) {
	topic, status := m.auth.UserIDFor(r)
	if status != 200 {
		log.Printf("Auth Error: %d", status)
		w.WriteHeader(status)
		return
	}
  conn, err := m.upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("Upgrade error: %v", err)
    return
  }
	conns := m.createConns(topic)
	id := conns.Add(conn)
	m.pubSub.Subscribe(topic)
	defer m.removeConn(id, conns)
	for {
    if _, _, err := conn.NextReader(); err != nil {
      conn.Close()
      break
    }
  }
}

func (m *Mux) deliver(topic string, message string) {
	// If you look at the source, NewPreparedMessage never actually
	// raises an error.
	pmsg, _ := websocket.NewPreparedMessage(websocket.TextMessage, []byte(message))
	conns := m.connsFor(topic)
	if conns == nil {
		return
	}
	conns.Each(func(conn *websocket.Conn)error {
		return conn.WritePreparedMessage(pmsg)
	})
}

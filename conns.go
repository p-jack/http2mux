package wsmux

import "log"
import "sync"
import "github.com/gorilla/websocket"

type conns struct {
  topic string
  index int
  conns map[int]*websocket.Conn
  mutex *sync.Mutex
}

func newConns(topic string) *conns {
  return &conns{
    topic: topic,
    index: 0,
    conns: map[int]*websocket.Conn{},
    mutex: &sync.Mutex{},
  }
}

func (c *conns) Add(conn *websocket.Conn) int {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  c.index++
  c.conns[c.index] = conn
  return c.index
}

func (c *conns) Remove(id int) int {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  delete(c.conns, id)
  return len(c.conns)
}

func (c *conns) list() []*websocket.Conn {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  // TODO: This is inefficient and nasty. :(
  result := make([]*websocket.Conn, 0, len(c.conns))
  for _, conn := range c.conns {
    result = append(result, conn)
  }
  return result
}

func (c *conns) Each(sink func(conn *websocket.Conn)error) {
  list := c.list()
  for _, conn := range list {
    go func(conn *websocket.Conn) {
      err := sink(conn)
      if err != nil {
        log.Printf("sink: %v", err)
        conn.Close();
      }
    }(conn)
  }
}

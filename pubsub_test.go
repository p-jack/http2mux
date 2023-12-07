package wsmux

import "sync"

type TestPubSub struct {
  mutex *sync.Mutex
  deliver func(string,string)
  subs map[string]int
}

func NewTestPubSub() *TestPubSub {
  return &TestPubSub{
    mutex:&sync.Mutex{},
    subs:map[string]int{},
  }
}

func (ps *TestPubSub) Subscribe(topic string) {
  ps.mutex.Lock()
  defer ps.mutex.Unlock()
  ps.subs[topic] = ps.subs[topic] + 1
}

func (ps *TestPubSub) Unsubscribe(topic string) {
  ps.mutex.Lock()
  defer ps.mutex.Unlock()
  if ps.subs[topic] != 1 {
    panic("assertion failed (Unsubscribe should only be called when no subs are left)")
  }
  delete(ps.subs, topic)
}

func (ps *TestPubSub) DeliverTo(f func(topic string, message string)) {
  ps.mutex.Lock()
  defer ps.mutex.Unlock()
  ps.deliver = f
}

func (ps *TestPubSub) Publish(topic string, message string) {
  ps.deliver(topic, message)
}

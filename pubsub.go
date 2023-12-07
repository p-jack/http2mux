package wsmux

type PubSub interface {
  Subscribe(topic string)
  Unsubscribe(topic string)
  DeliverTo(func(topic string, message string))
}

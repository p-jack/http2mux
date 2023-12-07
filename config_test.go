package wsmux

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewConfig(t *testing.T) {
  assert := assert.New(t)
  config := NewConfig()
  assert.Equal(config.Addr, ":8080")
  assert.Equal(config.Endpoint, "/ws")
}

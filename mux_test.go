package wsmux

import "net/http"
import "net/url"
import "testing"
import "github.com/stretchr/testify/require"
import "github.com/gorilla/websocket"

type test struct {
  auth *TestAuth
  pubSub *TestPubSub
  mux *Mux
}

func setUp() test {
  auth := NewTestAuth()
  pubSub := NewTestPubSub()
  mux := New(NewConfig(), auth, pubSub)
  mux.Start()
  return test{
    auth: auth,
    pubSub: pubSub,
    mux: mux,
  }
}

func (t test) tearDown() {
  t.mux.Stop()
}

func (t test) header(session string) http.Header {
  header := http.Header{}
  if session != "" {
    cookie := http.Cookie{
      Name: "session-id",
      Value: session,
    }
    header.Add("Cookie", cookie.String())
  }
  return header
}

func (t test) dial(session string) (*websocket.Conn, error) {
  header := t.header(session)
  u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
  conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
  return conn, err
}

func TestSimple(t *testing.T) {
  require := require.New(t)
  test := setUp()
  defer test.tearDown()
  test.auth.Add("session1", "userA")
  conn, err := test.dial("session1")
  defer conn.Close()
  require.NoError(err)
  test.pubSub.Publish("userA", "Hello, world!")
  _, message, err := conn.ReadMessage()
  require.NoError(err)
  require.Equal("Hello, world!", string(message))
}

func TestNoAuth(t *testing.T) {
  require := require.New(t)
  test := setUp()
  defer test.tearDown()
  _, err := test.dial("")
  require.Equal("websocket: bad handshake", err.Error())
}

func TestBadUpgrade(t *testing.T) {
  require := require.New(t)
  test := setUp()
  defer test.tearDown()
  test.auth.Add("session1", "userA")
  header := test.header("session1")
  client := &http.Client{}
  req, err := http.NewRequest("GET", "http://localhost:8080/ws", nil)
  require.NoError(err)
  req.Header = header
  resp, err := client.Do(req)
  require.NoError(err)
  require.Equal(400, resp.StatusCode)
}

func TestNoSubs(t *testing.T) {
  //require := require.New(t)
  test := setUp()
  defer test.tearDown()
  test.pubSub.Publish("userA", "Hello, world!")
}

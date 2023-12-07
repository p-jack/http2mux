package wsmux

import "sync"
import "net/http"

type TestAuth struct {
  mutex *sync.Mutex
  sessions map[string]string
}

func NewTestAuth() *TestAuth {
  return &TestAuth{
    mutex: &sync.Mutex{},
    sessions: map[string]string{},
  }
}

func (a *TestAuth) UserIDFor(request *http.Request) (string,int) {
  cookie, err := request.Cookie("session-id")
  if err != nil {
    return "", 401
  }
  session := cookie.Value
  a.mutex.Lock()
  defer a.mutex.Unlock()
  user := a.sessions[session]
  if user == "" {
    return "", 401
  } else {
    return user, 200
  }
}

func (a *TestAuth) Add(session string, user string) {
  a.mutex.Lock()
  defer a.mutex.Unlock()
  a.sessions[session] = user
}

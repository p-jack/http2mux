package wsmux

import "net/http"

type Auth interface {
  UserIDFor(request *http.Request) (string,int)
}

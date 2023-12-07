package wsmux

type Config struct {
	Addr string `env:"WSMUX_ADDR" envDefault:":8080"`
  Endpoint string `env:"WSMUX_ADDR" envDefault:"/ws"`
	UserCookie string `env:"WSMUX_USER_COOKIE" envDefault:"user-id"`
}

func NewConfig() Config {
  return Config{Addr:":8080", Endpoint:"/ws", UserCookie:"user-id"}
}

package conf

import (
	"flag"
	"goout/lib/http"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf conf
	Conf = &Config{}
)

// Config config
type Config struct {
	HTTPServer *ServerConfig
	HTTPClient *http.ClientConfig
}

// ServerConfig Http Servers conf.
type ServerConfig struct {
	Addr string
}

func init() {
	flag.StringVar(&confPath, "conf", "goout-example.toml", "config path")
}

// Init init conf
func Init() (err error) {
	if _, err = toml.DecodeFile(confPath, &Conf); err != nil {
		return
	}
	return
}

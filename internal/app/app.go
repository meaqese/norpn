package app

import (
	"github.com/meaqese/norpn/internal/transport/rest"
	"log"
	"net/http"
	"os"
)

type Application struct {
	config *Config
}

type Config struct {
	port string
}

func ConfigFromEnv() *Config {
	config := &Config{}
	config.port = os.Getenv("PORT")
	if config.port == "" {
		config.port = "8080"
	}

	return config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() {
	http.HandleFunc("/api/v1/calculate", rest.CalcHandler)

	addr := ":" + a.config.port

	log.Println("Server started on " + addr)
	http.ListenAndServe(addr, nil)
}

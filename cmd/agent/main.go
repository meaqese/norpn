package main

import agent "github.com/meaqese/norpn/internal/agent/app"

func main() {
	app := agent.New()
	app.Run()
}

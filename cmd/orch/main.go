package main

import (
	"github.com/meaqese/norpn/internal/orch/app"
)

func main() {
	application := orch.New()
	err := application.Run()
	if err != nil {
		panic(err)
	}
}

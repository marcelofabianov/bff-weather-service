package main

import (
	"github.com/marcelofabianov/weather-server/internal/di"
)

func main() {
	di.NewApp().Run()
}

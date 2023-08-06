package main

import (
	"github.com/TadaTeruki/peruservice-auth/api"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()

	// Start server
	e.Logger.Fatal(api.NewServer(e).Start())
}

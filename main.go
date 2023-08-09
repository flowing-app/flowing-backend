package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var port int

func init() {
	flag.IntVar(&port, "p", 8080, "port number to be listen")
	flag.Parse()
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handleOk)
	e.GET("/echo", handleEcho)
	e.POST("/runn", handleRunn)

	e.Logger.SetLevel(log.INFO)
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", port)))
}

func handleOk(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func handleEcho(c echo.Context) error {
	text := c.QueryParam("text")
	return c.String(http.StatusOK, text)
}

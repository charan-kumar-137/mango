package mango

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	MaxThreads      uint64 // TODO
	RouteMap        RouteMap
	InitialContext  context.Context
	shutDownHandler func(context.Context)
}

func CreateApp() App {
	return App{}
}

func (app *App) AddRoutes(routeMap RouteMap) {
	app.RouteMap = routeMap
	app.ValidateRoutes()
}

func (app *App) ValidateRoutes() {
	for route, config := range app.RouteMap {
		if config.HandlerFunc == nil {
			panic(fmt.Sprintf("No Handler Function for %s\n", route))
		}
	}
}

func (app *App) AddShutDownHandler(f func(context.Context)) {
	app.shutDownHandler = f
}

func (app *App) Run(host string, port string) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var addr string = fmt.Sprintf("%s:%s", host, port)
	log.Printf("Running on %s", addr)

	var server http.Server = http.Server{Addr: addr, Handler: app.RouteMap}
	go func() {
		sig := <-sigs
		log.Printf("Received Signal - %s; Shutting Down Server", sig)
		if app.shutDownHandler != nil {
			app.shutDownHandler(app.InitialContext)
		}
		if err := server.Shutdown(app.InitialContext); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
	}()
	var err error = server.ListenAndServe()

	log.Fatalf("Error Occurred while Running: %s", err.Error())

}

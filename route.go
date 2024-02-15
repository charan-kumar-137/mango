package mango

import (
	"context"
	"log"
	"net/http"
	"time"
)

type RouteConfig struct {
	AllowedMethods []string
	HandlerFunc    func(context.Context, Request) Response
}

type RouteMap map[string]RouteConfig

func (routeMap RouteMap) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var c context.Context = context.Background()

	var start time.Time = time.Now()
	var path string = req.URL.Path

	var response Response

	defer func() {
		if r := recover(); r != nil {
			log.Println("Internal Error: ", r)
			response = Response{Data: []byte("Internal Server Error"), Status: http.StatusInternalServerError}
		}

		response.Send(w)

		var end time.Time = time.Now()

		log.Printf("%s %d %s", path, response.Status, end.Sub(start))
	}()
	request := ParseRequest(req)
	if routeConfig, ok := routeMap[path]; ok {
		response = routeConfig.HandlerFunc(c, request)
	} else {
		response = Response{Data: []byte("Not Found"), Status: http.StatusNotFound}
	}
}

package main

import (
	"fmt"
	"gin-photo-storage/conf"
	"gin-photo-storage/routers"
	"gin-photo-storage/constant"
	"net/http"
)

func main() {
	// get the global router
	router := routers.Router

	// set up a http server
	server := http.Server{
		Addr: fmt.Sprintf(":%s", conf.ServerCfg.Get(constant.SERVER_PORT)),
		Handler: router,
		MaxHeaderBytes: 1 << 20,
	}

	// run the server
	server.ListenAndServeTLS("conf/server.crt", "conf/server.key")
	//server.ListenAndServe()
}

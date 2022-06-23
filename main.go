package main

import (
	"net/http"

	"github.com/junaozun/web_framework_demo/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	myRun(core)
}

func myRun(core *framework.Core) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}
	server.ListenAndServe()
}

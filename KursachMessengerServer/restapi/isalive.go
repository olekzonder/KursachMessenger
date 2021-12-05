package restapi

import (
	"KursachMessenger/api"
	"net/http"
	"os"
)

var IsAlive = func(w http.ResponseWriter,r *http.Request) {
	name := os.Getenv("server_name")
	api.Respond(w,api.MakeResponse(true,name))
	return
}

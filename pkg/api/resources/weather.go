package resources

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	router *httprouter.Router
}

func (res *Weather) FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! This is the first endpoint working!")
}

func (res *Weather) Register(router *httprouter.Router) {
	res.router = router
	res.router.GET("/", res.FirstEndpoint)
}

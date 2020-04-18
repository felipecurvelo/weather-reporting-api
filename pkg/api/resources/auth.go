package resources

import (
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	router *httprouter.Router
}

func (res *Auth) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (res *Auth) Register(router *httprouter.Router) {
	res.router = router
	res.router.POST("/auth/", res.Authorize)
}

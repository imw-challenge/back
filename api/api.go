package api

import (
	"crypto/subtle"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imw-challenge/back/db"
)

type API struct {
	router *mux.Router
	mdb    *db.MessageDB
}

func InitAPI(db *db.MessageDB) (*API, error) {
	a := &API{
		router: mux.NewRouter(),
		mdb:    db,
	}
	a.SetRoutes()
	return a, nil
}

func (a *API) SetRoutes() {
	a.PublicPost("/public/message", a.postMessageHandler()) // unauthenticated
	a.PrivatePut("/private/message", a.putMessageHandler())
	a.PrivateGet("/private/message", a.getMessageHandler())
	a.PrivateGet("/private/dump", a.getDumpHandler())
}

func (a *API) GetRouter() *mux.Router {
	return a.router
}

func (a *API) PublicGet(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("GET")
}

func (a *API) PrivateGet(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, defaultAuth(f)).Methods("GET")
}

func (a *API) PublicPost(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("POST")
}

func (a *API) PrivatePost(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, defaultAuth(f)).Methods("POST")
}

func (a *API) PublicPut(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("PUT")
}

func (a *API) PrivatePut(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, defaultAuth(f)).Methods("PUT")
}

func (a *API) PublicDelete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("DELETE")
}

func (a *API) PrivateDelete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, defaultAuth(f)).Methods("DELETE")
}

func defaultAuth(handler http.HandlerFunc) http.HandlerFunc {
	return basicAuth(handler, "admin", "back-challenge", "Please enter credentials:")
}

func basicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

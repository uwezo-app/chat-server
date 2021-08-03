package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/uwezo-app/chat-server/controller"
	"github.com/uwezo-app/chat-server/server"
)

func Handlers(hub *server.Hub) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)

	r.HandleFunc("/register", controller.CreatePsychologist).Methods(http.MethodPost)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
	r.HandleFunc("/login", controller.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
	r.HandleFunc("/reset", controller.ResetHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", controller.LogoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		server.ChatHandler(hub, w, r)
	})

	r.Use(mux.CORSMethodMiddleware(r))
	// Auth route
	// s := r.PathPrefix("/auth").Subrouter()
	// s.Use(auth.JwtVerify)
	// s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	// s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")

	//c := r.PathPrefix("/ws").Subrouter()
	//c.HandleFunc("/chat", controllers.ChatHandler)  /ws/chat/

	return r
}

// CommonMiddleware --Set basic headers
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

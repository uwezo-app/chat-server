package router

import (
	"github.com/gorilla/mux"
	controllers "github.com/uwezo-app/chat-server/controller"
	"net/http"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)

	r.HandleFunc("/register", controllers.CreatePsychologist).Methods("POST")
	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", controllers.LogoutHandler).Methods("POST")
	r.HandleFunc("/chat", controllers.ChatHandler)

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

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
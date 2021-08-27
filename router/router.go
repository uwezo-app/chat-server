package router

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/gorilla/mux"

	"github.com/uwezo-app/chat-server/controller"
	"github.com/uwezo-app/chat-server/server"
)

func Handlers(hub *server.Hub, dbase *gorm.DB) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.CreatePsychologist(dbase, w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginHandler(dbase, w, r)
	}).Methods(http.MethodPost)

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)

	r.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		controller.ResetHandler(dbase, w, r)
	}).Methods(http.MethodPost)

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		controller.LogoutHandler(dbase, w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/psychologists/profile/{Email}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controller.GetProfileHandler(dbase, w, r)
		} else if r.Method == http.MethodPost {
			controller.UpdateProfileHandler(dbase, w, r)
		} else if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}
	}).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

	/**
	=================
		METRICS
	=================
	*/
	r.HandleFunc("/psychologists", func(w http.ResponseWriter, r *http.Request) {
		controller.GetPsychologists(dbase, w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/psychologists/number", func(w http.ResponseWriter, r *http.Request) {
		controller.GetNumberofPsychologists(dbase, w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/patients/number", func(w http.ResponseWriter, r *http.Request) {
		controller.GetNumberofPatients(dbase, w, r)
	}).Methods(http.MethodGet)
	/**
	=================
		END METRICS
	=================
	*/

	/**
	=================
		ADMIN
	=================
	*/
	r.HandleFunc("/admin/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.AdminRegistrationHandler(dbase, w, r)
		}
	}).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.AdminLoginHandler(dbase, w, r)
		}
	}).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.AdminLogoutHandler(dbase, w, r)
		}
	}).Methods(http.MethodGet, http.MethodOptions)
	/**
	=================
		END ADMIN
	=================
	*/

	/**
	=================
	Patient
	=================
	*/
	r.HandleFunc("/patient/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.PatientLoginHandler(dbase, w, r)
		}
	}).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/patient/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.CreatePatient(dbase, w, r)
		}
	}).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/patient/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
		} else if r.Method == http.MethodPost {
			controller.PatientLogoutHandler(dbase, w, r)
		}
	}).Methods(http.MethodGet, http.MethodOptions)
	/**
	=================
	Patient
	=================
	*/

	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		server.ChatHandler(hub, dbase, w, r)
	})

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

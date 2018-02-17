package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App Abacus Backend
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize App
func (a *App) Initialize(dbhost, user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, dbhost, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		errlog.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run App
func (a *App) Run(addr string) {
	errlog.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.PathPrefix("/page/").Handler(http.StripPrefix("/page/", http.FileServer(http.Dir("frontend"))))
	a.Router.PathPrefix("/qr/").Handler(http.StripPrefix("/qr/", http.FileServer(http.Dir("qr"))))
	a.Router.HandleFunc("/participants/aid/{aid}", a.getParticipantByAID).Methods("GET")
	a.Router.HandleFunc("/participants/id/{id}", a.getParticipantByID).Methods("GET")
	a.Router.HandleFunc("/participants", a.createParticipant).Methods("POST")
	a.Router.HandleFunc("/participants/id/{id}", a.updateParticipant).Methods("POST")
	a.Router.HandleFunc("/allparticipants", a.getParticipants).Methods("GET")
	a.Router.HandleFunc("/allcolleges", a.getColleges).Methods("GET")
	a.Router.HandleFunc("/colleges/id/{id}", a.getCollege).Methods("GET")
	a.Router.HandleFunc("/checkin", a.checkIn).Methods("POST")
	a.Router.HandleFunc("/events/{id}/participants", a.getEventParticipants).Methods("GET")
	a.Router.HandleFunc("/participants/aid/{aid}/events", a.getEventsOfParticipant).Methods("GET")
	a.Router.HandleFunc("/checkout/{aid}", a.checkOut).Methods("GET")
	a.Router.HandleFunc("/certrequests", a.getCheckOuts).Methods("GET")
}

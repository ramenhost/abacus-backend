package main

import (
	"database/sql"
	"encoding/json"
	"helpers"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (a *App) getParticipantByAID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := participant{AID: vars["aid"]}
	if err := p.getParticipant(a.DB, "aid"); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getParticipantByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	p := participant{ID: id}
	if err := p.getParticipant(a.DB, "id"); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getParticipantByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := participant{Email: vars["email"]}
	if err := p.getParticipant(a.DB, "email"); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createParticipant(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p participant
	defer r.Body.Close()
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	if p.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	if err := p.createParticipant(a.DB); err != nil {
		if strings.Contains(err.Error(), "1062") {
			respondWithError(w, http.StatusForbidden, "User already Exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	go helpers.MailRegistrationLink(p.Email, p.ID)
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateParticipant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	p := participant{ID: id}
	if err := p.getParticipant(a.DB, "id"); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	if err := p.updateParticipant(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	go helpers.MailNewRegistration(p.Email, p.AID)
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getParticipants(w http.ResponseWriter, r *http.Request) {
	participants, err := getParticipants(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, participants)
}

func (a *App) getColleges(w http.ResponseWriter, r *http.Request) {
	colleges, err := getColleges(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, colleges)
}

func (a *App) getCollege(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	c := college{ID: id}
	if err := c.getCollege(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "College not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) checkIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EventID    int      `json:"eventNo"`
		Passphrase string   `json:"passphrase"`
		Entries    []string `json:"entries"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode body")
		return
	}
	var e = event{ID: req.EventID}
	if err := e.getPassphrase(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Event not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if e.Passphrase != req.Passphrase {
		respondWithError(w, http.StatusUnauthorized, "Passphrase doesn't match with Event")
		return
	}
	var accepted = []string{}
	for aid := range req.Entries {
		if err := e.checkIn(a.DB, req.Entries[aid]); err == nil {
			accepted = append(accepted, req.Entries[aid])
		}
	}
	var res struct {
		Pushed []string `json:"pushed"`
	}
	res.Pushed = accepted
	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) checkOut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := participant{AID: vars["aid"]}
	if err := p.checkOut(a.DB); err != nil {
		if strings.Contains(err.Error(), "1062") {
			respondWithError(w, http.StatusForbidden, "Certificate Requested Already")
			return
		} else if strings.Contains(err.Error(), "1452") {
			respondWithError(w, http.StatusForbidden, "User Not Found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Success string `json:"message"`
	}{"Certificate Requested"})
}

func (a *App) getEventParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	e := event{ID: id}
	var participants []participant
	if participants, err = e.getParticipants(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, participants)
}

func (a *App) getEventsOfParticipant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := participant{AID: vars["aid"]}
	var events []string
	var err error
	if events, err = p.getEvents(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, events)
}

func (a *App) getCheckOuts(w http.ResponseWriter, r *http.Request) {
	participants, err := getCheckOuts(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, participants)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	errlog.Println(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

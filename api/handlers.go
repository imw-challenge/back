package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/imw-challenge/back/types"
)

func (a *API) postMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m types.Message
		err := decoder.Decode(&m)
		if err != nil {
			badRequestHandler(w, "postMessage", err)
			return
		}
		if len(m.ID) == 0 || len(m.Text) == 0 {
			badRequestHandler(w, "postMessage", errors.New("No ID or Text in request"))
			return
		}
		err = a.mdb.InsertMessage(&m)
		if err != nil {
			internalErrorHandler(w, "postMessage", err)
			return
		}
		w.Write([]byte{})
	}
}

func (a *API) putMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m types.Message
		err := decoder.Decode(&m)
		if err != nil {
			badRequestHandler(w, "putMessage", err)
			return
		}
		if len(m.ID) == 0 || len(m.Text) == 0 {
			badRequestHandler(w, "putMessage", errors.New("No ID or Text in request"))
			return
		}
		message, err := a.mdb.FetchByID(m.ID)
		if err != nil {
			notFoundHandler(w, m.ID, "putMessage - fetchyByID")
			return
		}
		message.Text = m.Text
		err = a.mdb.InsertMessage(message)
		if err != nil {
			internalErrorHandler(w, "putMessage", err)
			return
		}
		w.Write([]byte{})
	}
}

func (a *API) getMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m types.Message
		err := decoder.Decode(&m)
		if err != nil {
			badRequestHandler(w, "getMessage", err)
			return
		}
		if len(m.ID) == 0 {
			badRequestHandler(w, "getMessages", errors.New("No ID in request"))
			return
		}
		message, err := a.mdb.FetchByID(m.ID)
		if err != nil {
			notFoundHandler(w, m.ID, "getMessage - fetchyByID")
			return
		}
		messageJSON, err := json.MarshalIndent(message, "", "    ")
		if err != nil {
			internalErrorHandler(w, "getMessage", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(messageJSON)
	}
}

func (a *API) getDumpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := a.mdb.FetchAntiChrono()
		if err != nil {
			internalErrorHandler(w, "getDump", err)
			return
		}
		messagesJSON, err := json.MarshalIndent(messages, "", "    ")
		if err != nil {
			internalErrorHandler(w, "getDump", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(messagesJSON)
	}
}

func internalErrorHandler(w http.ResponseWriter, handlerID string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("Internal error in %s: %s", handlerID, err)
}

func badRequestHandler(w http.ResponseWriter, handlerID string, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Printf("Bad Request in %s: %s", handlerID, err)
}

func notFoundHandler(w http.ResponseWriter, resourceID string, handlerID string) {
	w.WriteHeader(http.StatusNotFound)
	log.Printf("Resource %s not found in %s", resourceID, handlerID)
}

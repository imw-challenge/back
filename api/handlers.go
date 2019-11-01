package api

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"

	"github.com/imw-challenge/back/types"
)

// postMessageHandler handles a post message request
// it checks that there is a well-formed body, containing at least an ID and text
// and inserts to the DB
func (a *API) postMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			badRequestHandler(w, "postMessage", errors.New("Request had no body"))
			return
		}
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

// putMessageHandler handles a put message request
// it checks that the request has a well-formed body containing an ID and text
// if a message with this ID exists, it updates the text
// otherwise, it returns 404
func (a *API) putMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body == nil {
			badRequestHandler(w, "putMessage", errors.New("Request had no body"))
			return
		}

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

// getMessageHandler handles a get message request
// it checks that there is a well-formed request body containing an ID
// it returns 404 if this message is not in the db, otherwise returning
// the message as pretty-printed JSON
func (a *API) getMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			badRequestHandler(w, "getMessage", errors.New("Request had no body"))
			return
		}

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

// getDumpHandler handles a get dump request
// it fetches all of the messages in reverse chronoligcal order,
// and returns them as a pretty-printed JSON array
func (a *API) getDumpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//this fetches all messages, starting with the latest
		messages, err := a.mdb.FetchSortedByTime(0, math.MaxInt64, false)
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
}

func notFoundHandler(w http.ResponseWriter, resourceID string, handlerID string) {
	w.WriteHeader(http.StatusNotFound)
}

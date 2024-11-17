package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AkhilKJames/rssaggregator/internal/auth"
	"github.com/AkhilKJames/rssaggregator/internal/database"
)

type authhandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authhandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			log.Println("ERROR: Authentication ", err)
			respondWithError(w, 403, fmt.Sprintf("error authenticating: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			log.Println("ERROR: No user with that api key: ", err)
			respondWithError(w, 400, fmt.Sprintf("error user dosent exist: %v", err))
			return
		}

		handler(w, r, user)
	}
}

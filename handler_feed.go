package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AkhilKJames/rssaggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("ERROR: Cant parse JSON: ", err)
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		log.Println("ERROR: Couldnt update table feed: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt create feed: %v", err))
		return
	}
	respondWithJson(w, 201, databaseFeedtoFeed(feed))
}

func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		log.Println("ERROR: Couldnt get feeds: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt get feeds: %v", err))
		return
	}
	respondWithJson(w, 201, databaseFeedstoFeeds(feeds))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AkhilKJames/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("ERROR: Cant parse JSON: ", err)
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		log.Println("ERROR: Couldnt create feed follow: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt create feed follow: %v", err))
		return
	}
	respondWithJson(w, 201, databaseFeedFollowtoFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Println("ERROR: Couldnt get feed follows: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt get feed follows: %v", err))
		return
	}
	respondWithJson(w, 200, databaseFeedFollowstoFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDString := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		log.Println("ERROR: invalid uuid: ", err)
		respondWithError(w, 400, fmt.Sprintf("error invalid feed follow id: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		log.Println("ERROR: failed to delete record: ", err)
		respondWithError(w, 500, fmt.Sprintf("error failed to delete feed follow: %v", err))
		return
	}
	respondWithJson(w, 200, struct{}{})
}

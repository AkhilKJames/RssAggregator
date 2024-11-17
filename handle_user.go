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

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("ERROR: Cant parse JSON: ", err)
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		log.Println("ERROR: Couldnt update table user: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt create user: %v", err))
		return
	}
	respondWithJson(w, 201, databaseUsertoUser(user))
}

func (apiCfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, databaseUsertoUser(user))
}

func (apiCfg *apiConfig) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		log.Println("ERROR: Couldnt get posts followed by user: ", err)
		respondWithError(w, 400, fmt.Sprintf("error couldnt get posts followed by user:: %v", err))
		return
	}
	respondWithJson(w, 200, databasePoststoPosts(posts))
}

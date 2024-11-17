package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AkhilKJames/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portstring := os.Getenv("PORT")
	if portstring == "" {
		log.Fatal("PORT not defined in the environement")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL URL not defined in the environement")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("ERROR: Cannot connect to database: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	// scraping feed
	db := database.New(conn)
	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1router := chi.NewRouter()
	v1router.Get("/healthz", handleReadiness)
	v1router.Get("/err", handleError)
	v1router.Post("/users", apiCfg.handleCreateUser)
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUserByApiKey))

	v1router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1router.Get("/feeds", apiCfg.handleGetFeeds)

	v1router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))

	v1router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollows))
	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

	router.Mount("/v1", v1router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portstring,
	}

	log.Printf("Server starting on Port %v", portstring)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}

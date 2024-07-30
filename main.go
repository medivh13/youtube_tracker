package main

import (
	"log"
	"net/http"
	"os"
	"youtube_tracker/src/handler"
	"youtube_tracker/src/infra/youtube"
	repository "youtube_tracker/src/repository/video"
	useCase "youtube_tracker/src/usecase/video"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	dsn := os.Getenv("POSTGRES_DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	videoRepo := repository.NewVideoRepository(db)
	youtubeClient, err := youtube.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Could not create YouTube client: %v", err)
	}

	videoService := useCase.NewVideoService(videoRepo, youtubeClient)
	handler := handler.NewHandler(videoService)

	r := chi.NewRouter()
	r.Get("/search", handler.SearchVideos)
	r.Get("/videos", handler.GetVideos)
	r.Get("/analitic", handler.GetPopularityScores)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}

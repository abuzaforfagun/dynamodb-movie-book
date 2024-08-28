package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type MovieAddedHandler struct {
	movieService services.MovieService
	genreService services.GenreService
}

func NewMovieAddedHandler(movieService *services.MovieService, genreService *services.GenreService) *MovieAddedHandler {
	return &MovieAddedHandler{
		movieService: *movieService,
		genreService: *genreService,
	}
}

func (h *MovieAddedHandler) HandleMessage(msg amqp.Delivery) {
	var payload events.MovieCreated
	log.Printf("Processing message [MessageId=%s]", payload.MessageId)

	err := json.Unmarshal(msg.Body, &payload)

	if err != nil {
		log.Println("Unable to unmarshal", err)
		return
	}

	if payload.MovieId == "" {
		log.Println("ERROR: MovieId should not be empty.")
		return
	}

	movie, err := h.movieService.GetInfo(payload.MovieId)

	if err != nil || movie == nil {
		log.Printf("ERROR: Invalid [MovieId=%s]\n", payload.MovieId)
		return
	}

	err = h.genreService.AddMovieToGenres(payload.MovieId, movie.Title, movie.ReleaseYear, movie.Genres)

	if err != nil {
		log.Println("ERROR: Unable to populate movies under genres", err)
	}
	log.Printf("Message processing completed [MessageId=%s]", payload.MessageId)
}

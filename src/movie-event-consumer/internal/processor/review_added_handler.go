package processor

import (
	"encoding/json"
	"log"
	"math"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type ReviewAddedHandler struct {
	movieService  services.MovieService
	reviewService services.ReviewService
}

func NewReviewAddedHandler(
	movieService *services.MovieService,
	reviewService *services.ReviewService) ReviewAddedHandler {
	return ReviewAddedHandler{
		movieService:  *movieService,
		reviewService: *reviewService,
	}
}

func (h *ReviewAddedHandler) HandleMessage(msg amqp.Delivery) {
	var payload events.ReviewAdded
	log.Printf("Processing message [MessageId=%s]", payload.MessageId)

	err := json.Unmarshal(msg.Body, &payload)

	if err != nil {
		log.Println("Unable to unmarshal", err)
		msg.Nack(false, false)
		return
	}

	if payload.MovieId == "" || payload.UserId == "" {
		log.Println("ERROR: MovieId should not be empty.")
		msg.Nack(false, false)
		return
	}

	movie, err := h.movieService.GetInfo(payload.MovieId)

	if err != nil || movie == nil {
		log.Printf("ERROR: Invalid [MovieId=%s]\n", payload.MovieId)
		msg.Nack(false, false)
		return
	}

	reviews, err := h.reviewService.GetReviews(payload.MovieId)
	if err != nil {
		log.Println("ERROR: Unable to get reviews")
		msg.Nack(false, false)
		return
	}

	totalScore := 0.0
	for _, review := range reviews {
		totalScore += review.Score
	}

	avgScore := totalScore / float64(len(reviews))
	roundedAvgScore := math.Round(avgScore*100) / 100
	err = h.movieService.UpdateMovieScore(payload.MovieId, roundedAvgScore)
	if err != nil {
		log.Println("ERROR: Unable to update movie score")
		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
	log.Printf("Processing completed [MessageId=%s]", payload.MessageId)
}

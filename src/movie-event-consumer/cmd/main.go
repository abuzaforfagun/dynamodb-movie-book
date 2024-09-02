package main

import (
	"log"
	"os"
	"strconv"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	enviornment := os.Getenv("ENVOIRNMENT")

	if enviornment != "production" {
		initializers.LoadEnvVariables()
	}

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")
	dynamodbUrl := os.Getenv("DYNAMODB_URL")

	dbConfig := dynamodb_connector.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
		Url:          dynamodbUrl,
	}

	dbConnector, err := dynamodb_connector.New(&dbConfig)

	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	reviewAddedExchangeName := os.Getenv("EXCHANGE_NAME_REVIEW_ADDED")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")
	reviewAddedQueueName := os.Getenv("REVIEW_ADDED_QUEUE")

	numberOfTopRatedMoviesStr := os.Getenv("NUMBER_OF_TOP_MOVIES")
	numberOfTopRatedMovies, err := strconv.Atoi(numberOfTopRatedMoviesStr)
	if err != nil {
		log.Fatal("Faild to initialize the consumer", err)
	}

	movieScoreUpdatedQueueName := os.Getenv("MOVIE_SCORE_UPDATED_QUEUE")

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	rmq, err := rabbitmq.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatal("Unable to connect to RabbitMQ", err)
	}
	defer rmq.Close()

	publisher, err := rabbitmq.NewPublisher(rabbitMqUri)
	if err != nil {
		log.Fatal("Unable to create publisher", err)
	}
	defer publisher.Close()
	movieScoreUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_SCORE_UPDATED")

	genreService := services.NewGenreService(dbConnector.Client, awsTableName)
	movieService := services.NewMovieService(dbConnector.Client, publisher, awsTableName, movieScoreUpdatedExchangeName, numberOfTopRatedMovies)
	reviewService := services.NewReviewService(dbConnector.Client, awsTableName)
	dlxService := services.NewDlxService(dbConnector.Client, awsTableName)

	moviedAddedHandler := processor.NewMovieAddedHandler(&movieService, &genreService)
	reviewAddedHandler := processor.NewReviewAddedHandler(&movieService, &reviewService)
	movieScoreUpdatedHandler := processor.NewMovieScoreUpdatedHandler(&movieService)
	dlxMovieAddedHandler := processor.NewDlxMovieAddedHandler(dlxService)
	dlxReviewAddedHandler := processor.NewDlxMovieAddedHandler(dlxService)

	dlxExchangeName := os.Getenv("DLX")
	dlxQueueName := os.Getenv("DLX")
	rmq.DeclareTopicExchanges([]string{dlxExchangeName})

	rmq.RegisterQueueExchange(dlxQueueName, dlxExchangeName, movieAddedQueueName, nil, dlxMovieAddedHandler.HandleMessage)
	movieAddedTable := amqp.Table{
		"x-message-ttl":             int32(10000),
		"x-dead-letter-exchange":    dlxExchangeName,     // The DLX exchange
		"x-dead-letter-routing-key": movieAddedQueueName, // Routing key for DLX
	}
	rmq.RegisterQueueExchange(movieAddedQueueName, movieAddedExchangeName, "", movieAddedTable, moviedAddedHandler.HandleMessage)

	rmq.RegisterQueueExchange(dlxQueueName, dlxExchangeName, reviewAddedQueueName, nil, dlxReviewAddedHandler.HandleMessage)
	reviewAddedTable := amqp.Table{
		"x-message-ttl":             int32(10000),
		"x-dead-letter-exchange":    dlxExchangeName,      // The DLX exchange
		"x-dead-letter-routing-key": reviewAddedQueueName, // Routing key for DLX
	}
	rmq.RegisterQueueExchange(reviewAddedQueueName, reviewAddedExchangeName, "", reviewAddedTable, reviewAddedHandler.HandleMessage)

	rmq.DeclareDirectExchanges([]string{movieScoreUpdatedExchangeName})
	rmq.RegisterQueueExchange(movieScoreUpdatedQueueName, movieScoreUpdatedExchangeName, "", nil, movieScoreUpdatedHandler.HandleMessage)

	log.Println("Ready to process events...")
	select {}
}

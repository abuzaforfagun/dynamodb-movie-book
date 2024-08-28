package grpc_services

import (
	"context"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
)

type ActorsService struct {
	repository repositories.ActorRepository
	actorpb.UnimplementedActorsServiceServer
}

func NewActorService(repository repositories.ActorRepository) *ActorsService {
	return &ActorsService{
		repository: repository,
	}
}

func (s *ActorsService) GetActorBasicInfo(ctx context.Context, request *actorpb.GetActorBasicInforRequestModel) (*actorpb.GetActorBasicInforResponseModel, error) {
	if len(request.ActorIds) == 0 {
		return nil, nil
	}
	dbResult, err := s.repository.Get(request.ActorIds)
	if err != nil {
		log.Println("ERROR: unable to get actor basic info", err)
		return nil, err
	}
	actors := []*actorpb.ActorBasicInfo{}

	for _, actor := range *dbResult {

		model := actorpb.ActorBasicInfo{
			Id:   actor.Id,
			Name: actor.Name,
		}

		actors = append(actors, &model)
	}

	result := &actorpb.GetActorBasicInforResponseModel{
		Actors: actors,
	}
	return result, nil
}

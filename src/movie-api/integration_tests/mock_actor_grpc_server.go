package integration_tests

import (
	"context"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
)

type MockActorGrpcServer struct {
	actorpb.UnimplementedActorsServiceServer
}

func (s *MockActorGrpcServer) GetActorBasicInfo(ctx context.Context, req *actorpb.GetActorBasicInforRequestModel) (*actorpb.GetActorBasicInforResponseModel, error) {

	var result []*actorpb.ActorBasicInfo
	for _, actorId := range req.ActorIds {
		if actorId != ValidActor1Id && actorId != ValidActor2Id {
			return &actorpb.GetActorBasicInforResponseModel{HasError: true}, nil
		}
		result = append(result, &actorpb.ActorBasicInfo{Id: actorId, Name: "Jhon"})
	}

	return &actorpb.GetActorBasicInforResponseModel{
		Actors: result,
	}, nil
}

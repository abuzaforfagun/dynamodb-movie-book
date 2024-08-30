package integration_tests

import (
	"context"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
)

type MockUserGrpcServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *MockUserGrpcServer) GetUserBasicInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {

	if req.UserId == ValidUserId {
		return &userpb.GetUserInfoResponse{
			Id:    req.UserId,
			Name:  "Jack",
			Email: "jack@jack.com",
		}, nil
	}

	return &userpb.GetUserInfoResponse{
		HasError: true,
	}, nil
}

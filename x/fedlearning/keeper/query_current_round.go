package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	"flmainchain/x/fedlearning/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetCurrentRound(ctx context.Context, req *types.QueryGetCurrentRoundRequest) (*types.QueryGetCurrentRoundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.CurrentRound.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetCurrentRoundResponse{CurrentRound: val}, nil
}

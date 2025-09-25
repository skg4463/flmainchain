package keeper

import (
	"context"
	"errors"

	"flmainchain/x/fedlearning/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListRound(ctx context.Context, req *types.QueryAllRoundRequest) (*types.QueryAllRoundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rounds, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Round,
		req.Pagination,
		func(_ uint64, value types.Round) (types.Round, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRoundResponse{Round: rounds, Pagination: pageRes}, nil
}

func (q queryServer) GetRound(ctx context.Context, req *types.QueryGetRoundRequest) (*types.QueryGetRoundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Round.Get(ctx, req.RoundId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetRoundResponse{Round: val}, nil
}

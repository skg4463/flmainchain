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

func (q queryServer) ListSubmittedScore(ctx context.Context, req *types.QueryAllSubmittedScoreRequest) (*types.QueryAllSubmittedScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	submittedScores, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.SubmittedScore,
		req.Pagination,
		func(_ string, value types.SubmittedScore) (types.SubmittedScore, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSubmittedScoreResponse{SubmittedScore: submittedScores, Pagination: pageRes}, nil
}

func (q queryServer) GetSubmittedScore(ctx context.Context, req *types.QueryGetSubmittedScoreRequest) (*types.QueryGetSubmittedScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.SubmittedScore.Get(ctx, req.ScoreId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetSubmittedScoreResponse{SubmittedScore: val}, nil
}

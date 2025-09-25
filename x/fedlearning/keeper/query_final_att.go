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

func (q queryServer) ListFinalAtt(ctx context.Context, req *types.QueryAllFinalAttRequest) (*types.QueryAllFinalAttResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	finalAtts, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.FinalAtt,
		req.Pagination,
		func(_ uint64, value types.FinalAtt) (types.FinalAtt, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllFinalAttResponse{FinalAtt: finalAtts, Pagination: pageRes}, nil
}

func (q queryServer) GetFinalAtt(ctx context.Context, req *types.QueryGetFinalAttRequest) (*types.QueryGetFinalAttResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.FinalAtt.Get(ctx, req.RoundId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetFinalAttResponse{FinalAtt: val}, nil
}

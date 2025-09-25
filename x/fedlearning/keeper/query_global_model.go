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

func (q queryServer) ListGlobalModel(ctx context.Context, req *types.QueryAllGlobalModelRequest) (*types.QueryAllGlobalModelResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	globalModels, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.GlobalModel,
		req.Pagination,
		func(_ uint64, value types.GlobalModel) (types.GlobalModel, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGlobalModelResponse{GlobalModel: globalModels, Pagination: pageRes}, nil
}

func (q queryServer) GetGlobalModel(ctx context.Context, req *types.QueryGetGlobalModelRequest) (*types.QueryGetGlobalModelResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.GlobalModel.Get(ctx, req.RoundId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetGlobalModelResponse{GlobalModel: val}, nil
}

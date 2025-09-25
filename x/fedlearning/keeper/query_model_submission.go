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

func (q queryServer) ListModelSubmission(ctx context.Context, req *types.QueryAllModelSubmissionRequest) (*types.QueryAllModelSubmissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	modelSubmissions, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.ModelSubmission,
		req.Pagination,
		func(_ string, value types.ModelSubmission) (types.ModelSubmission, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllModelSubmissionResponse{ModelSubmission: modelSubmissions, Pagination: pageRes}, nil
}

func (q queryServer) GetModelSubmission(ctx context.Context, req *types.QueryGetModelSubmissionRequest) (*types.QueryGetModelSubmissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.ModelSubmission.Get(ctx, req.OriginalHash)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetModelSubmissionResponse{ModelSubmission: val}, nil
}

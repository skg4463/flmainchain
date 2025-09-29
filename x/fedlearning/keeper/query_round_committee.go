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

func (q queryServer) ListRoundCommittee(ctx context.Context, req *types.QueryAllRoundCommitteeRequest) (*types.QueryAllRoundCommitteeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	roundCommittees, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.RoundCommittee,
		req.Pagination,
		func(_ uint64, value types.RoundCommittee) (types.RoundCommittee, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRoundCommitteeResponse{RoundCommittee: roundCommittees, Pagination: pageRes}, nil
}

func (q queryServer) GetRoundCommittee(ctx context.Context, req *types.QueryGetRoundCommitteeRequest) (*types.QueryGetRoundCommitteeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.RoundCommittee.Get(ctx, req.RoundId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetRoundCommitteeResponse{RoundCommittee: val}, nil
}

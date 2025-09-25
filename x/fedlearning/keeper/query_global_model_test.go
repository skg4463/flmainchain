package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"flmainchain/x/fedlearning/keeper"
	"flmainchain/x/fedlearning/types"
)

func createNGlobalModel(keeper keeper.Keeper, ctx context.Context, n int) []types.GlobalModel {
	items := make([]types.GlobalModel, n)
	for i := range items {
		items[i].RoundId = uint64(i)
		items[i].OriginalHash = strconv.Itoa(i)
		_ = keeper.GlobalModel.Set(ctx, items[i].RoundId, items[i])
	}
	return items
}

func TestGlobalModelQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNGlobalModel(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetGlobalModelRequest
		response *types.QueryGetGlobalModelResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetGlobalModelRequest{
				RoundId: msgs[0].RoundId,
			},
			response: &types.QueryGetGlobalModelResponse{GlobalModel: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetGlobalModelRequest{
				RoundId: msgs[1].RoundId,
			},
			response: &types.QueryGetGlobalModelResponse{GlobalModel: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetGlobalModelRequest{
				RoundId: 100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetGlobalModel(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestGlobalModelQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNGlobalModel(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGlobalModelRequest {
		return &types.QueryAllGlobalModelRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListGlobalModel(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GlobalModel), step)
			require.Subset(t, msgs, resp.GlobalModel)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListGlobalModel(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GlobalModel), step)
			require.Subset(t, msgs, resp.GlobalModel)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListGlobalModel(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.GlobalModel)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListGlobalModel(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

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

func createNRoundCommittee(keeper keeper.Keeper, ctx context.Context, n int) []types.RoundCommittee {
	items := make([]types.RoundCommittee, n)
	for i := range items {
		items[i].RoundId = uint64(i)
		items[i].Members = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		_ = keeper.RoundCommittee.Set(ctx, items[i].RoundId, items[i])
	}
	return items
}

func TestRoundCommitteeQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRoundCommittee(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRoundCommitteeRequest
		response *types.QueryGetRoundCommitteeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRoundCommitteeRequest{
				RoundId: msgs[0].RoundId,
			},
			response: &types.QueryGetRoundCommitteeResponse{RoundCommittee: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRoundCommitteeRequest{
				RoundId: msgs[1].RoundId,
			},
			response: &types.QueryGetRoundCommitteeResponse{RoundCommittee: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRoundCommitteeRequest{
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
			response, err := qs.GetRoundCommittee(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestRoundCommitteeQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRoundCommittee(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRoundCommitteeRequest {
		return &types.QueryAllRoundCommitteeRequest{
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
			resp, err := qs.ListRoundCommittee(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoundCommittee), step)
			require.Subset(t, msgs, resp.RoundCommittee)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListRoundCommittee(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoundCommittee), step)
			require.Subset(t, msgs, resp.RoundCommittee)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListRoundCommittee(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.RoundCommittee)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListRoundCommittee(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

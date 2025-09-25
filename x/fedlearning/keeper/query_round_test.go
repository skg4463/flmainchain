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

func createNRound(keeper keeper.Keeper, ctx context.Context, n int) []types.Round {
	items := make([]types.Round, n)
	for i := range items {
		items[i].RoundId = uint64(i)
		items[i].Status = strconv.Itoa(i)
		items[i].RequiredLNodes = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		items[i].SubmittedLNodes = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		items[i].RequiredCNodes = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		items[i].SubmittedCNodes = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		_ = keeper.Round.Set(ctx, items[i].RoundId, items[i])
	}
	return items
}

func TestRoundQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRound(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRoundRequest
		response *types.QueryGetRoundResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRoundRequest{
				RoundId: msgs[0].RoundId,
			},
			response: &types.QueryGetRoundResponse{Round: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRoundRequest{
				RoundId: msgs[1].RoundId,
			},
			response: &types.QueryGetRoundResponse{Round: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRoundRequest{
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
			response, err := qs.GetRound(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestRoundQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRound(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRoundRequest {
		return &types.QueryAllRoundRequest{
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
			resp, err := qs.ListRound(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Round), step)
			require.Subset(t, msgs, resp.Round)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListRound(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Round), step)
			require.Subset(t, msgs, resp.Round)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListRound(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Round)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListRound(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

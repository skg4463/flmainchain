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

func createNSubmittedScore(keeper keeper.Keeper, ctx context.Context, n int) []types.SubmittedScore {
	items := make([]types.SubmittedScore, n)
	for i := range items {
		items[i].ScoreId = strconv.Itoa(i)
		items[i].LnodeAddresses = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		items[i].Scores = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		_ = keeper.SubmittedScore.Set(ctx, items[i].ScoreId, items[i])
	}
	return items
}

func TestSubmittedScoreQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNSubmittedScore(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetSubmittedScoreRequest
		response *types.QueryGetSubmittedScoreResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetSubmittedScoreRequest{
				ScoreId: msgs[0].ScoreId,
			},
			response: &types.QueryGetSubmittedScoreResponse{SubmittedScore: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetSubmittedScoreRequest{
				ScoreId: msgs[1].ScoreId,
			},
			response: &types.QueryGetSubmittedScoreResponse{SubmittedScore: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSubmittedScoreRequest{
				ScoreId: strconv.Itoa(100000),
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
			response, err := qs.GetSubmittedScore(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestSubmittedScoreQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNSubmittedScore(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSubmittedScoreRequest {
		return &types.QueryAllSubmittedScoreRequest{
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
			resp, err := qs.ListSubmittedScore(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.SubmittedScore), step)
			require.Subset(t, msgs, resp.SubmittedScore)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListSubmittedScore(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.SubmittedScore), step)
			require.Subset(t, msgs, resp.SubmittedScore)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListSubmittedScore(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.SubmittedScore)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListSubmittedScore(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

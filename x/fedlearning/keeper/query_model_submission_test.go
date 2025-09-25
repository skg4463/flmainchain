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

func createNModelSubmission(keeper keeper.Keeper, ctx context.Context, n int) []types.ModelSubmission {
	items := make([]types.ModelSubmission, n)
	for i := range items {
		items[i].OriginalHash = strconv.Itoa(i)
		items[i].Tag = strconv.Itoa(i)
		items[i].Submitter = strconv.Itoa(i)
		_ = keeper.ModelSubmission.Set(ctx, items[i].OriginalHash, items[i])
	}
	return items
}

func TestModelSubmissionQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNModelSubmission(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetModelSubmissionRequest
		response *types.QueryGetModelSubmissionResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetModelSubmissionRequest{
				OriginalHash: msgs[0].OriginalHash,
			},
			response: &types.QueryGetModelSubmissionResponse{ModelSubmission: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetModelSubmissionRequest{
				OriginalHash: msgs[1].OriginalHash,
			},
			response: &types.QueryGetModelSubmissionResponse{ModelSubmission: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetModelSubmissionRequest{
				OriginalHash: strconv.Itoa(100000),
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
			response, err := qs.GetModelSubmission(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestModelSubmissionQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNModelSubmission(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllModelSubmissionRequest {
		return &types.QueryAllModelSubmissionRequest{
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
			resp, err := qs.ListModelSubmission(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ModelSubmission), step)
			require.Subset(t, msgs, resp.ModelSubmission)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListModelSubmission(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ModelSubmission), step)
			require.Subset(t, msgs, resp.ModelSubmission)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListModelSubmission(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.ModelSubmission)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListModelSubmission(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

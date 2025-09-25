package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"flmainchain/x/fedlearning/keeper"
	"flmainchain/x/fedlearning/types"
)

func TestCurrentRoundMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)
	expected := &types.MsgCreateCurrentRound{Creator: creator}
	_, err = srv.CreateCurrentRound(f.ctx, expected)
	require.NoError(t, err)
	rst, err := f.keeper.CurrentRound.Get(f.ctx)
	require.Nil(t, err)
	require.Equal(t, expected.Creator, rst.Creator)
}

func TestCurrentRoundMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	expected := &types.MsgCreateCurrentRound{Creator: creator}
	_, err = srv.CreateCurrentRound(f.ctx, expected)
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateCurrentRound
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgUpdateCurrentRound{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgUpdateCurrentRound{Creator: unauthorizedAddr},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "completed",
			request: &types.MsgUpdateCurrentRound{Creator: creator},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateCurrentRound(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, err := f.keeper.CurrentRound.Get(f.ctx)
				require.Nil(t, err)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestCurrentRoundMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateCurrentRound(f.ctx, &types.MsgCreateCurrentRound{Creator: creator})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteCurrentRound
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgDeleteCurrentRound{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgDeleteCurrentRound{Creator: unauthorizedAddr},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "completed",
			request: &types.MsgDeleteCurrentRound{Creator: creator},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteCurrentRound(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				found, err := f.keeper.CurrentRound.Has(f.ctx)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}

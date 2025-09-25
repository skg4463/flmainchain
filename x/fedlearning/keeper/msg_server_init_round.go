package keeper

import (
	"context"

	"flmainchain/x/fedlearning/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) InitRound(ctx context.Context, msg *types.MsgInitRound) (*types.MsgInitRoundResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgInitRoundResponse{}, nil
}

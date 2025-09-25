package keeper

import (
	"context"

	"flmainchain/x/fedlearning/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitWeight(ctx context.Context, msg *types.MsgSubmitWeight) (*types.MsgSubmitWeightResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSubmitWeightResponse{}, nil
}

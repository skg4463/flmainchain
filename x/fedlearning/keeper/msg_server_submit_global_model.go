package keeper

import (
	"context"

	"flmainchain/x/fedlearning/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitGlobalModel(ctx context.Context, msg *types.MsgSubmitGlobalModel) (*types.MsgSubmitGlobalModelResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSubmitGlobalModelResponse{}, nil
}

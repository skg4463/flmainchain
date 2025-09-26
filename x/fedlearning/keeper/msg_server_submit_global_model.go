package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) SubmitGlobalModel(goCtx context.Context, msg *types.MsgSubmitGlobalModel) (*types.MsgSubmitGlobalModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// k.RoundCommittee -> k.Keeper.RoundCommittee, msg.RoundID -> msg.RoundId
	committee, err := k.Keeper.RoundCommittee.Get(ctx, msg.RoundId)
	if err != nil { return nil, errorsmod.Wrapf(types.ErrRoundNotFound, "committee for round %d not found: %s", msg.RoundId, err.Error()) }
	if len(committee.Members) == 0 || msg.Creator != committee.Members[0] {
		return nil, errorsmod.Wrapf(types.ErrNotCommitteeLeader, "address %s is not the committee leader", msg.Creator)
	}

	// k.GlobalModel -> k.Keeper.GlobalModel, msg.RoundID -> msg.RoundId
	err = k.Keeper.GlobalModel.Set(ctx, msg.RoundId, types.GlobalModel{
		OriginalHash: msg.OriginalHash,
	})
	if err != nil { return nil, err }

	return &types.MsgSubmitGlobalModelResponse{}, nil
}
package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) InitRound(goCtx context.Context, msg *types.MsgInitRound) (*types.MsgInitRoundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	_, err := k.Keeper.CurrentRound.Get(ctx)
	if err == nil {
		return nil, types.ErrAlreadyInitialized
	}

	k.Keeper.CurrentRound.Set(ctx, types.CurrentRound{RoundId: 1})
	k.Keeper.Round.Set(ctx, 1, types.Round{
		RoundId:         1,
		Status:          "WeightSubmissionOpen",
		RequiredLNodes:  msg.InitialLNodes,
		SubmittedLNodes: []string{},
		RequiredCNodes:  msg.InitialCNodes,
		SubmittedCNodes: []string{},
	})
	// Roundcommittee -> RoundCommittee 로 수정
	k.Keeper.RoundCommittee.Set(ctx, 1, types.RoundCommittee{RoundId: 1, Members: msg.InitialCNodes})

	return &types.MsgInitRoundResponse{}, nil
}
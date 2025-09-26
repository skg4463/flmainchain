package keeper

import (
	"context"
	// "errorsmod"는 사용하지 않으므로 삭제
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) InitRound(goCtx context.Context, msg *types.MsgInitRound) (*types.MsgInitRoundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// k.CurrentRound -> k.Keeper.CurrentRound
	_, err := k.Keeper.CurrentRound.Get(ctx)
	// Go의 if문은 boolean 값을 기대함. err != nil 로 수정.
	if err == nil {
		return nil, types.ErrAlreadyInitialized
	}

	// k.CurrentRound -> k.Keeper.CurrentRound, RoundID -> RoundId
	k.Keeper.CurrentRound.Set(ctx, types.CurrentRound{RoundId: 1})
	k.Keeper.Round.Set(ctx, 1, types.Round{
		RoundId:         1,
		Status:          "WeightSubmissionOpen",
		RequiredLNodes:  msg.InitialMembers,
		SubmittedLNodes: []string{},
		RequiredCNodes:  msg.InitialMembers,
		SubmittedCNodes: []string{},
	})
	k.Keeper.RoundCommittee.Set(ctx, 1, types.RoundCommittee{Members: msg.InitialMembers})

	return &types.MsgInitRoundResponse{}, nil
}
package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) SubmitGlobalModel(goCtx context.Context, msg *types.MsgSubmitGlobalModel) (*types.MsgSubmitGlobalModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Roundcommittee -> RoundCommittee 로 수정
	committee, err := k.Keeper.RoundCommittee.Get(ctx, msg.RoundId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrRoundNotFound, "committee for round %d not found: %s", msg.RoundId, err.Error())
	}
	if len(committee.Members) == 0 || msg.Creator != committee.Members[0] {
		return nil, errorsmod.Wrapf(types.ErrNotCommitteeLeader, "address %s is not the committee leader", msg.Creator)
	}

	err = k.Keeper.GlobalModel.Set(ctx, msg.RoundId, types.GlobalModel{
		OriginalHash: msg.OriginalHash,
	})
	if err != nil { return nil, err }

	k.Keeper.ElectNextCommittee(ctx)
	//if err != nil {
	//	// 만약 위원회 선출에 실패하면, 전체 트랜잭션을 롤백하기 위해 에러를 반환.
	//	return nil, errorsmod.Wrap(err, "failed to elect next committee")
	//}

	return &types.MsgSubmitGlobalModelResponse{}, nil
}
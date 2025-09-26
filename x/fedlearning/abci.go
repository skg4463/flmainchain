package fedlearning

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/keeper"
)

// BeginBlocker. 매 블록 시작 시 실행.
func BeginBlocker(ctx context.Context, k keeper.Keeper) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.AdvanceRoundState(sdkCtx)
}

// EndBlocker. 매 블록 끝에서 실행.
func EndBlocker(ctx context.Context, k keeper.Keeper) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// 3N-1 블록 (높이가 2, 5, 8, ...)의 끝에서 ATT 자동 집계.
	if (height > 1) && (height%3) == 2 {
		k.AggregateScoresAndCreateATT(sdkCtx)
	}
	// 3N 블록 (높이가 3, 6, 9, ...)의 끝에서 다음 라운드 위원회 선출.
	if (height > 1) && (height%3) == 0 {
		k.ElectNextCommittee(sdkCtx)
	}
}
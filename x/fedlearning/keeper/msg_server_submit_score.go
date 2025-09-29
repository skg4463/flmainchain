package keeper

import (
	"context"
	"fmt"
	"slices"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) SubmitScore(goCtx context.Context, msg *types.MsgSubmitScore) (*types.MsgSubmitScoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	round, err := k.Keeper.Round.Get(ctx, msg.RoundId)
	if err != nil { return nil, errorsmod.Wrapf(types.ErrRoundNotFound, "round %d: %s", msg.RoundId, err.Error()) }
	if round.Status != "ScoreSubmissionOpen" { return nil, errorsmod.Wrapf(types.ErrInvalidRoundStatus, "expected ScoreSubmissionOpen, got %s", round.Status) }
	if !slices.Contains(round.RequiredCNodes, msg.Creator) { return nil, errorsmod.Wrapf(types.ErrNotAParticipant, "address %s is not a required C-node", msg.Creator) }
	if slices.Contains(round.SubmittedCNodes, msg.Creator) { return nil, errorsmod.Wrapf(types.ErrAlreadySubmitted, "address %s has already submitted scores", msg.Creator) }
	if len(msg.LnodeAddresses) != len(msg.Scores) { return nil, errorsmod.Wrapf(types.ErrInvalidData, "address count (%d) and score count (%d) do not match", len(msg.LnodeAddresses), len(msg.Scores)) }

	id := fmt.Sprintf("%d-%s", msg.RoundId, msg.Creator)
	// Creator 필드 삭제
	err = k.Keeper.SubmittedScore.Set(ctx, id, types.SubmittedScore{
		LnodeAddresses: msg.LnodeAddresses,
		Scores:         msg.Scores,
	})
	if err != nil { return nil, err }

	round.SubmittedCNodes = append(round.SubmittedCNodes, msg.Creator)
	err = k.Keeper.Round.Set(ctx, round.RoundId, round)
	if err != nil { return nil, err }

	return &types.MsgSubmitScoreResponse{}, nil
}
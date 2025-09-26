package keeper

import (
	"context"
	"slices"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k msgServer) SubmitWeight(goCtx context.Context, msg *types.MsgSubmitWeight) (*types.MsgSubmitWeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// k.Round -> k.Keeper.Round, msg.RoundID -> msg.RoundId
	round, err := k.Keeper.Round.Get(ctx, msg.RoundId)
	if err != nil { return nil, errorsmod.Wrapf(types.ErrRoundNotFound, "round %d: %s", msg.RoundId, err.Error()) }
	if round.Status != "WeightSubmissionOpen" { return nil, errorsmod.Wrapf(types.ErrInvalidRoundStatus, "expected WeightSubmissionOpen, got %s", round.Status) }
	if !slices.Contains(round.RequiredLNodes, msg.Creator) { return nil, errorsmod.Wrapf(types.ErrNotAParticipant, "address %s is not a required L-node", msg.Creator) }
	if slices.Contains(round.SubmittedLNodes, msg.Creator) { return nil, errorsmod.Wrapf(types.ErrAlreadySubmitted, "address %s has already submitted", msg.Creator) }

	// k.ModelSubmission -> k.Keeper.ModelSubmission
	err = k.Keeper.ModelSubmission.Set(ctx, msg.OriginalHash, types.ModelSubmission{
		Submitter: msg.Creator,
		Tag:       msg.Tag,
	})
	if err != nil { return nil, err }

	round.SubmittedLNodes = append(round.SubmittedLNodes, msg.Creator)
	// k.Round -> k.Keeper.Round, round.RoundID -> round.RoundId
	err = k.Keeper.Round.Set(ctx, round.RoundId, round)
	if err != nil { return nil, err }

	return &types.MsgSubmitWeightResponse{}, nil
}
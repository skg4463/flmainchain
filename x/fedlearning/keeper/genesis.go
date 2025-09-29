package keeper

import (
	"context"
	"errors"

	"flmainchain/x/fedlearning/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}
	if genState.CurrentRound != nil {
		if err := k.CurrentRound.Set(ctx, *genState.CurrentRound); err != nil {
			return err
		}
	}
	for _, elem := range genState.RoundMap {
		if err := k.Round.Set(ctx, elem.RoundId, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.ModelSubmissionMap {
		if err := k.ModelSubmission.Set(ctx, elem.OriginalHash, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.SubmittedScoreMap {
		if err := k.SubmittedScore.Set(ctx, elem.ScoreId, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.FinalAttMap {
		if err := k.FinalAtt.Set(ctx, elem.RoundId, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.GlobalModelMap {
		if err := k.GlobalModel.Set(ctx, elem.RoundId, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.RoundCommitteeMap {
		if err := k.RoundCommittee.Set(ctx, elem.RoundId, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	currentRound, err := k.CurrentRound.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	genesis.CurrentRound = &currentRound
	if err := k.Round.Walk(ctx, nil, func(_ uint64, val types.Round) (stop bool, err error) {
		genesis.RoundMap = append(genesis.RoundMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.ModelSubmission.Walk(ctx, nil, func(_ string, val types.ModelSubmission) (stop bool, err error) {
		genesis.ModelSubmissionMap = append(genesis.ModelSubmissionMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.SubmittedScore.Walk(ctx, nil, func(_ string, val types.SubmittedScore) (stop bool, err error) {
		genesis.SubmittedScoreMap = append(genesis.SubmittedScoreMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.FinalAtt.Walk(ctx, nil, func(_ uint64, val types.FinalAtt) (stop bool, err error) {
		genesis.FinalAttMap = append(genesis.FinalAttMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.GlobalModel.Walk(ctx, nil, func(_ uint64, val types.GlobalModel) (stop bool, err error) {
		genesis.GlobalModelMap = append(genesis.GlobalModelMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.RoundCommittee.Walk(ctx, nil, func(_ uint64, val types.RoundCommittee) (stop bool, err error) {
		genesis.RoundCommitteeMap = append(genesis.RoundCommitteeMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}

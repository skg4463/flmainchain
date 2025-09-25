package keeper_test

import (
	"testing"

	"flmainchain/x/fedlearning/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:       types.DefaultParams(),
		PortId:       types.PortID,
		CurrentRound: &types.CurrentRound{RoundId: 93}, RoundMap: []types.Round{{RoundId: 0}, {RoundId: 1}}, ModelSubmissionMap: []types.ModelSubmission{{OriginalHash: "0"}, {OriginalHash: "1"}}, SubmittedScoreMap: []types.SubmittedScore{{ScoreId: "0"}, {ScoreId: "1"}}, FinalAttMap: []types.FinalAtt{{RoundId: 0}, {RoundId: 1}}, GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)
	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.CurrentRound, got.CurrentRound)
	require.EqualExportedValues(t, genesisState.RoundMap, got.RoundMap)
	require.EqualExportedValues(t, genesisState.ModelSubmissionMap, got.ModelSubmissionMap)
	require.EqualExportedValues(t, genesisState.SubmittedScoreMap, got.SubmittedScoreMap)
	require.EqualExportedValues(t, genesisState.FinalAttMap, got.FinalAttMap)
	require.EqualExportedValues(t, genesisState.GlobalModelMap, got.GlobalModelMap)

}

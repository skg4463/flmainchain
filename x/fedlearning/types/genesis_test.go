package types_test

import (
	"testing"

	"flmainchain/x/fedlearning/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId:       types.PortID,
				CurrentRound: &types.CurrentRound{RoundId: 53}, RoundMap: []types.Round{{RoundId: 0}, {RoundId: 1}}, ModelSubmissionMap: []types.ModelSubmission{{OriginalHash: "0"}, {OriginalHash: "1"}}, SubmittedScoreMap: []types.SubmittedScore{{ScoreId: "0"}, {ScoreId: "1"}}, FinalAttMap: []types.FinalAtt{{RoundId: 0}, {RoundId: 1}}, GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}},
			valid: true,
		}, {
			desc: "duplicated round",
			genState: &types.GenesisState{
				RoundMap: []types.Round{
					{
						RoundId: 0,
					},
					{
						RoundId: 0,
					},
				},
				ModelSubmissionMap: []types.ModelSubmission{{OriginalHash: "0"}, {OriginalHash: "1"}}, SubmittedScoreMap: []types.SubmittedScore{{ScoreId: "0"}, {ScoreId: "1"}}, FinalAttMap: []types.FinalAtt{{RoundId: 0}, {RoundId: 1}}, GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}},
			valid: false,
		}, {
			desc: "duplicated modelSubmission",
			genState: &types.GenesisState{
				ModelSubmissionMap: []types.ModelSubmission{
					{
						OriginalHash: "0",
					},
					{
						OriginalHash: "0",
					},
				},
				SubmittedScoreMap: []types.SubmittedScore{{ScoreId: "0"}, {ScoreId: "1"}}, FinalAttMap: []types.FinalAtt{{RoundId: 0}, {RoundId: 1}}, GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}},
			valid: false,
		}, {
			desc: "duplicated submittedScore",
			genState: &types.GenesisState{
				SubmittedScoreMap: []types.SubmittedScore{
					{
						ScoreId: "0",
					},
					{
						ScoreId: "0",
					},
				},
				FinalAttMap: []types.FinalAtt{{RoundId: 0}, {RoundId: 1}}, GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}},
			valid: false,
		}, {
			desc: "duplicated finalAtt",
			genState: &types.GenesisState{
				FinalAttMap: []types.FinalAtt{
					{
						RoundId: 0,
					},
					{
						RoundId: 0,
					},
				},
				GlobalModelMap: []types.GlobalModel{{RoundId: 0}, {RoundId: 1}}},
			valid: false,
		}, {
			desc: "duplicated globalModel",
			genState: &types.GenesisState{
				GlobalModelMap: []types.GlobalModel{
					{
						RoundId: 0,
					},
					{
						RoundId: 0,
					},
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

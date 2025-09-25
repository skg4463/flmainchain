package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		PortId: PortID, CurrentRound: nil, RoundMap: []Round{}, ModelSubmissionMap: []ModelSubmission{}, SubmittedScoreMap: []SubmittedScore{}, FinalAttMap: []FinalAtt{}, GlobalModelMap: []GlobalModel{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	roundIndexMap := make(map[string]struct{})

	for _, elem := range gs.RoundMap {
		index := fmt.Sprint(elem.RoundId)
		if _, ok := roundIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for round")
		}
		roundIndexMap[index] = struct{}{}
	}
	modelSubmissionIndexMap := make(map[string]struct{})

	for _, elem := range gs.ModelSubmissionMap {
		index := fmt.Sprint(elem.OriginalHash)
		if _, ok := modelSubmissionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for modelSubmission")
		}
		modelSubmissionIndexMap[index] = struct{}{}
	}
	submittedScoreIndexMap := make(map[string]struct{})

	for _, elem := range gs.SubmittedScoreMap {
		index := fmt.Sprint(elem.ScoreId)
		if _, ok := submittedScoreIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for submittedScore")
		}
		submittedScoreIndexMap[index] = struct{}{}
	}
	finalAttIndexMap := make(map[string]struct{})

	for _, elem := range gs.FinalAttMap {
		index := fmt.Sprint(elem.RoundId)
		if _, ok := finalAttIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for finalAtt")
		}
		finalAttIndexMap[index] = struct{}{}
	}
	globalModelIndexMap := make(map[string]struct{})

	for _, elem := range gs.GlobalModelMap {
		index := fmt.Sprint(elem.RoundId)
		if _, ok := globalModelIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for globalModel")
		}
		globalModelIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

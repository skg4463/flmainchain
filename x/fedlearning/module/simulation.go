package fedlearning

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	fedlearningsimulation "flmainchain/x/fedlearning/simulation"
	"flmainchain/x/fedlearning/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	fedlearningGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&fedlearningGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateCurrentRound          = "op_weight_msg_fedlearning"
		defaultWeightMsgCreateCurrentRound int = 100
	)

	var weightMsgCreateCurrentRound int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCurrentRound, &weightMsgCreateCurrentRound, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCurrentRound = defaultWeightMsgCreateCurrentRound
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCurrentRound,
		fedlearningsimulation.SimulateMsgCreateCurrentRound(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateCurrentRound          = "op_weight_msg_fedlearning"
		defaultWeightMsgUpdateCurrentRound int = 100
	)

	var weightMsgUpdateCurrentRound int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCurrentRound, &weightMsgUpdateCurrentRound, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCurrentRound = defaultWeightMsgUpdateCurrentRound
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCurrentRound,
		fedlearningsimulation.SimulateMsgUpdateCurrentRound(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteCurrentRound          = "op_weight_msg_fedlearning"
		defaultWeightMsgDeleteCurrentRound int = 100
	)

	var weightMsgDeleteCurrentRound int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCurrentRound, &weightMsgDeleteCurrentRound, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCurrentRound = defaultWeightMsgDeleteCurrentRound
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCurrentRound,
		fedlearningsimulation.SimulateMsgDeleteCurrentRound(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgInitRound          = "op_weight_msg_fedlearning"
		defaultWeightMsgInitRound int = 100
	)

	var weightMsgInitRound int
	simState.AppParams.GetOrGenerate(opWeightMsgInitRound, &weightMsgInitRound, nil,
		func(_ *rand.Rand) {
			weightMsgInitRound = defaultWeightMsgInitRound
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInitRound,
		fedlearningsimulation.SimulateMsgInitRound(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSubmitWeight          = "op_weight_msg_fedlearning"
		defaultWeightMsgSubmitWeight int = 100
	)

	var weightMsgSubmitWeight int
	simState.AppParams.GetOrGenerate(opWeightMsgSubmitWeight, &weightMsgSubmitWeight, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitWeight = defaultWeightMsgSubmitWeight
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitWeight,
		fedlearningsimulation.SimulateMsgSubmitWeight(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSubmitScore          = "op_weight_msg_fedlearning"
		defaultWeightMsgSubmitScore int = 100
	)

	var weightMsgSubmitScore int
	simState.AppParams.GetOrGenerate(opWeightMsgSubmitScore, &weightMsgSubmitScore, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitScore = defaultWeightMsgSubmitScore
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitScore,
		fedlearningsimulation.SimulateMsgSubmitScore(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSubmitGlobalModel          = "op_weight_msg_fedlearning"
		defaultWeightMsgSubmitGlobalModel int = 100
	)

	var weightMsgSubmitGlobalModel int
	simState.AppParams.GetOrGenerate(opWeightMsgSubmitGlobalModel, &weightMsgSubmitGlobalModel, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitGlobalModel = defaultWeightMsgSubmitGlobalModel
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitGlobalModel,
		fedlearningsimulation.SimulateMsgSubmitGlobalModel(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

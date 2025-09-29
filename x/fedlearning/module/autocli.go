package fedlearning

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"flmainchain/x/fedlearning/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "GetCurrentRound",
					Use:       "get-current-round",
					Short:     "Gets a CurrentRound",
					Alias:     []string{"show-current-round"},
				},
				{
					RpcMethod: "ListRound",
					Use:       "list-round",
					Short:     "List all Round",
				},
				{
					RpcMethod:      "GetRound",
					Use:            "get-round [id]",
					Short:          "Gets a Round",
					Alias:          []string{"show-round"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				{
					RpcMethod: "ListModelSubmission",
					Use:       "list-model-submission",
					Short:     "List all ModelSubmission",
				},
				{
					RpcMethod:      "GetModelSubmission",
					Use:            "get-model-submission [id]",
					Short:          "Gets a ModelSubmission",
					Alias:          []string{"show-model-submission"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "original_hash"}},
				},
				{
					RpcMethod: "ListSubmittedScore",
					Use:       "list-submitted-score",
					Short:     "List all SubmittedScore",
				},
				{
					RpcMethod:      "GetSubmittedScore",
					Use:            "get-submitted-score [id]",
					Short:          "Gets a SubmittedScore",
					Alias:          []string{"show-submitted-score"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "score_id"}},
				},
				{
					RpcMethod: "ListFinalAtt",
					Use:       "list-final-att",
					Short:     "List all FinalATT",
				},
				{
					RpcMethod:      "GetFinalAtt",
					Use:            "get-final-att [id]",
					Short:          "Gets a FinalATT",
					Alias:          []string{"show-final-att"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				{
					RpcMethod: "ListGlobalModel",
					Use:       "list-global-model",
					Short:     "List all GlobalModel",
				},
				{
					RpcMethod:      "GetGlobalModel",
					Use:            "get-global-model [id]",
					Short:          "Gets a GlobalModel",
					Alias:          []string{"show-global-model"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				{
					RpcMethod: "ListRoundCommittee",
					Use:       "list-round-committee",
					Short:     "List all RoundCommittee",
				},
				{
					RpcMethod:      "GetRoundCommittee",
					Use:            "get-round-committee [id]",
					Short:          "Gets a RoundCommittee",
					Alias:          []string{"show-round-committee"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateCurrentRound",
					Use:            "create-current-round [round-id]",
					Short:          "Create CurrentRound",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				{
					RpcMethod:      "UpdateCurrentRound",
					Use:            "update-current-round [round-id]",
					Short:          "Update CurrentRound",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}},
				},
				{
					RpcMethod: "DeleteCurrentRound",
					Use:       "delete-current-round",
					Short:     "Delete CurrentRound",
				},
				{
					RpcMethod:      "InitRound",
					Use:            "init-round [initial-members]",
					Short:          "Send a init-round tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "initial_members", Varargs: true}},
				},
				{
					RpcMethod:      "SubmitWeight",
					Use:            "submit-weight [round-id] [original-hash] [tag]",
					Short:          "Send a submit-weight tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}, {ProtoField: "original_hash"}, {ProtoField: "tag"}},
				},
				{
					RpcMethod:      "SubmitScore",
					Use:            "submit-score [round-id] [lnode-addresses] [scores]",
					Short:          "Send a submit-score tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}, {ProtoField: "lnode_addresses"}, {ProtoField: "scores", Varargs: true}},
				},
				{
					RpcMethod:      "SubmitGlobalModel",
					Use:            "submit-global-model [round-id] [original-hash]",
					Short:          "Send a submit-global-model tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "round_id"}, {ProtoField: "original_hash"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

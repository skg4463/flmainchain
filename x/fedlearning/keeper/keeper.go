package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"flmainchain/x/fedlearning/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Port collections.Item[string]

	ibcKeeperFn func() *ibckeeper.Keeper

	bankKeeper      types.BankKeeper
	stakingKeeper   types.StakingKeeper
	CurrentRound    collections.Item[types.CurrentRound]
	Round           collections.Map[uint64, types.Round]
	ModelSubmission collections.Map[string, types.ModelSubmission]
	SubmittedScore  collections.Map[string, types.SubmittedScore]
	FinalAtt        collections.Map[uint64, types.FinalAtt]
	GlobalModel     collections.Map[uint64, types.GlobalModel]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	ibcKeeperFn func() *ibckeeper.Keeper,

	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		ibcKeeperFn:   ibcKeeperFn,
		Port:          collections.NewItem(sb, types.PortKey, "port", collections.StringValue),
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		CurrentRound:  collections.NewItem(sb, types.CurrentRoundKey, "currentRound", codec.CollValue[types.CurrentRound](cdc)), Round: collections.NewMap(sb, types.RoundKey, "round", collections.Uint64Key, codec.CollValue[types.Round](cdc)), ModelSubmission: collections.NewMap(sb, types.ModelSubmissionKey, "modelSubmission", collections.StringKey, codec.CollValue[types.ModelSubmission](cdc)), SubmittedScore: collections.NewMap(sb, types.SubmittedScoreKey, "submittedScore", collections.StringKey, codec.CollValue[types.SubmittedScore](cdc)), FinalAtt: collections.NewMap(sb, types.FinalAttKey, "finalAtt", collections.Uint64Key, codec.CollValue[types.FinalAtt](cdc)), GlobalModel: collections.NewMap(sb, types.GlobalModelKey, "globalModel", collections.Uint64Key, codec.CollValue[types.GlobalModel](cdc))}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

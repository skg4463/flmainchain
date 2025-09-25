package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitGlobalModel{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitScore{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitWeight{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInitRound{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCurrentRound{},
		&MsgUpdateCurrentRound{},
		&MsgDeleteCurrentRound{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}

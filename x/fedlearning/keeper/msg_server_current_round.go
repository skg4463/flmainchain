package keeper

import (
	"context"
	"fmt"

	"flmainchain/x/fedlearning/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCurrentRound(ctx context.Context, msg *types.MsgCreateCurrentRound) (*types.MsgCreateCurrentRoundResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	found, err := k.CurrentRound.Has(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "already set")
	}

	var currentRound = types.CurrentRound{
		Creator: msg.Creator,
		RoundId: msg.RoundId,
	}

	if err := k.CurrentRound.Set(
		ctx,
		currentRound,
	); err != nil {
		return nil, err
	}

	return &types.MsgCreateCurrentRoundResponse{}, nil
}

func (k msgServer) UpdateCurrentRound(ctx context.Context, msg *types.MsgUpdateCurrentRound) (*types.MsgUpdateCurrentRoundResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value exists
	valFound, err := k.CurrentRound.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var currentRound = types.CurrentRound{
		Creator: msg.Creator,
		RoundId: msg.RoundId,
	}

	if err := k.CurrentRound.Set(ctx, currentRound); err != nil {
		return nil, err
	}

	return &types.MsgUpdateCurrentRoundResponse{}, nil
}

func (k msgServer) DeleteCurrentRound(ctx context.Context, msg *types.MsgDeleteCurrentRound) (*types.MsgDeleteCurrentRoundResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value exists
	val, err := k.CurrentRound.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.CurrentRound.Remove(ctx); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgDeleteCurrentRoundResponse{}, nil
}

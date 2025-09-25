package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/fedlearning module sentinel errors
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrRoundNotFound      = errors.Register(ModuleName, 1101, "round not found")
	ErrInvalidRoundStatus = errors.Register(ModuleName, 1102, "invalid round status for this action")
	ErrNotAParticipant    = errors.Register(ModuleName, 1103, "address is not a participant for this action")
	ErrAlreadySubmitted   = errors.Register(ModuleName, 1104, "address has already submitted for this round")
	ErrInvalidData        = errors.Register(ModuleName, 1105, "invalid data provided")
    ErrNotCommitteeLeader = errors.Register(ModuleName, 1106, "only committee leader can submit global model")
	ErrAlreadyInitialized = errors.Register(ModuleName, 1107, "round already initialized")
	ErrInvalidPacketTimeout = errors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errors.Register(ModuleName, 1501, "invalid version")
)

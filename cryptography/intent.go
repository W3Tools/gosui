package cryptography

import (
	"bytes"
)

// AppID is the application identifier for the intent.
type AppID = uint8

const (
	// Sui is the AppID for Sui.
	Sui AppID = iota
)

// IntentVersion is the version of the intent, which is used to manage changes in the intent structure over time.
type IntentVersion = uint8

const (
	// V0 is the IntentVersion for the initial version of Sui.
	V0 IntentVersion = iota
)

// IntentScope is the scope of the intent, indicating what type of data it applies to.
type IntentScope = uint8

const (
	// TransactionData is the IntentScope for transaction data.
	TransactionData IntentScope = iota
	// TransactionEffects is the IntentScope for transaction effects.
	TransactionEffects
	// CheckpointSummary is the IntentScope for checkpoint summaries.
	CheckpointSummary
	// PersonalMessage is the IntentScope for personal messages.
	PersonalMessage
)

// Intent defines the structure of an intent, which includes the scope, version, and application identifier.
type Intent = []uint8

// IntentWithScope creates an intent with the specified scope, using the default version and application identifier.
func IntentWithScope(scope IntentScope) Intent {
	return []uint8{scope, V0, Sui}
}

// MessageWithIntent combines the intent with the provided message, returning a new byte slice that includes both.
func MessageWithIntent(scope IntentScope, message []byte) []byte {
	intent := IntentWithScope(scope)
	intentMessage := new(bytes.Buffer)
	intentMessage.Write(intent)
	intentMessage.Write(message)

	return intentMessage.Bytes()
}

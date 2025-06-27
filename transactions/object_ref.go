package transactions

import (
	"fmt"
	"strconv"

	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/types"
)

// ObjectStringRef defines a string-based reference to a Sui object.
type ObjectStringRef struct {
	ObjectID string `json:"objectId"`
	Version  string `json:"version"`
	Digest   string `json:"digest"`
}

// ToObjectRef converts an ObjectStringRef to a sui_types.ObjectRef.
func (ref ObjectStringRef) ToObjectRef() (*sui_types.ObjectRef, error) {
	objectID, err := sui_types.NewObjectIdFromHex(ref.ObjectID)
	if err != nil {
		return nil, fmt.Errorf("can not create object id from hex [%s]: %v", ref.ObjectID, err)
	}
	digest, err := sui_types.NewDigest(ref.Digest)
	if err != nil {
		return nil, fmt.Errorf("can not create digest [%s]: %v", ref.Digest, err)
	}
	version, err := strconv.ParseUint(ref.Version, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can not parse version [%s] to uint64: %v", ref.Version, err)
	}

	return &sui_types.ObjectRef{ObjectId: *objectID, Version: version, Digest: *digest}, nil
}

func coinStructToObjectRef(coin types.CoinStruct) (*sui_types.ObjectRef, error) {
	return ObjectStringRef{ObjectID: coin.CoinObjectID, Version: coin.Version, Digest: coin.Digest}.ToObjectRef()
}

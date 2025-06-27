package types

import (
	"encoding/json"
	"errors"
)

// StakeObject is an interface that defines a Sui stake object.
type StakeObject interface {
	isStakeObject()
}

// StakeObjectPending defines a pending Sui stake object.
type StakeObjectPending struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiID       string `json:"stakedSuiId"`
	Status            string `json:"status"`
}

// StakeObjectActive defines an active Sui stake object.
type StakeObjectActive struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiID       string `json:"stakedSuiId"`
	EstimatedReward   string `json:"estimatedReward"`
	Status            string `json:"status"`
}

// StakeObjectUnstaked defines an unstaked Sui stake object.
type StakeObjectUnstaked struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiID       string `json:"stakedSuiId"`
	Status            string `json:"status"`
}

// isStakeObject implements the StakeObject interface for StakeObjectPending.
func (StakeObjectPending) isStakeObject() {}

// isStakeObject implements the StakeObject interface for StakeObjectActive.
func (StakeObjectActive) isStakeObject() {}

// isStakeObject implements the StakeObject interface for StakeObjectUnstaked.
func (StakeObjectUnstaked) isStakeObject() {}

// StakeObjectWrapper defines a wrapper for StakeObject to support custom JSON marshaling and unmarshaling.
type StakeObjectWrapper struct {
	StakeObject
}

// UnmarshalJSON decodes JSON data into a StakeObjectWrapper.
func (w *StakeObjectWrapper) UnmarshalJSON(data []byte) error {
	type Status struct {
		Status string `json:"status"`
	}
	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	switch status.Status {
	case "Pending":
		var so StakeObjectPending
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	case "Active":
		var so StakeObjectActive
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	case "Unstaked":
		var so StakeObjectUnstaked
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	default:
		return errors.New("unknown StakeObject type")
	}

	return nil
}

// MarshalJSON encodes a StakeObjectWrapper into JSON.
func (w *StakeObjectWrapper) MarshalJSON() ([]byte, error) {
	switch obj := w.StakeObject.(type) {
	case StakeObjectPending:
		return json.Marshal(StakeObjectPending{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiID:       obj.StakedSuiID,
		})
	case StakeObjectActive:
		return json.Marshal(StakeObjectActive{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiID:       obj.StakedSuiID,
			EstimatedReward:   obj.EstimatedReward,
		})
	case StakeObjectUnstaked:
		return json.Marshal(StakeObjectUnstaked{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiID:       obj.StakedSuiID,
		})
	default:
		return nil, errors.New("unknown StakeObject type")
	}
}

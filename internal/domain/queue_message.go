package domain

import (
	"errors"
)

// QueueMessage interface.
type QueueMessage interface {

	// GetMessageType - returns message type.
	GetMessageType() MessageType

	// Validate - validates message type.
	Validate() error
}

// PrepareGunsMessage - prepare guns event.
type PrepareGunsMessage struct {
}

// GetMessageType - returns PrepareGunsMessageType.
func (pgm *PrepareGunsMessage) GetMessageType() MessageType {
	return PrepareGunsMessageType
}

// Validate - validates message type.
func (pgm *PrepareGunsMessage) Validate() error {
	if !pgm.GetMessageType().IsValid() {
		return errors.New("Message type is not valid")
	}

	return nil
}

// ShootMessage - shoot event.
type ShootMessage struct {
	ShooterName string `json:"shooterName"`
	Damage      int32  `json:"damage"`
}

// GetMessageType - returns ShootMessageType.
func (sm *ShootMessage) GetMessageType() MessageType {
	return ShootMessageType
}

// Validate - validates message type.
func (sm *ShootMessage) Validate() error {
	if !sm.GetMessageType().IsValid() {
		return errors.New("Message type is not valid")
	}
	return nil
}

// WinnerMessage - winner event message.
type WinnerMessage struct {
	Message string `json:"message"`
}

// GetMessageType - returns WinnerMessageType.
func (wm *WinnerMessage) GetMessageType() MessageType {
	return WinnerMessageType
}

// Validate - validates message type.
func (wm *WinnerMessage) Validate() error {
	if !wm.GetMessageType().IsValid() {
		return errors.New("Message type is not valid")
	}
	return nil
}

var (
	_ QueueMessage = &PrepareGunsMessage{}
	_ QueueMessage = &ShootMessage{}
	_ QueueMessage = &WinnerMessage{}
)

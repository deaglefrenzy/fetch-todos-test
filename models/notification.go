package models

import "time"

type Notification struct {
	ID          string    `firestore:"-" json:"id"`
	User        User      `firestore:"user" json:"user"`
	Type        string    `firestore:"type" json:"type"`
	Message     string    `firestore:"message" json:"message"`
	ReferenceID string    `firestore:"referenceID" json:"reference_id"`
	Read        bool      `firestore:"read" json:"read"`
	CreatedAt   time.Time `firestore:"createdAt" json:"created_at"`
}

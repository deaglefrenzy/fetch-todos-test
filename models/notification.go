package models

import "time"

type Notification struct {
	ID        string      `firestore:"-" json:"id"`
	User      User        `firestore:"user" json:"user"`
	Message   string      `firestore:"message" json:"message"`
	Reference interface{} `firestore:"reference,omitempty" json:"reference,omitempty"`
	Read      bool        `firestore:"read" json:"read"`
	CreatedAt time.Time   `firestore:"createdAt" json:"created_at"`
}

type Reference struct {
	Type string      `firestore:"type"`
	Data interface{} `firestore:"data,omitempty"`
}

package notifications

import (
	"context"
	"learnfirestore/models"
	"time"

	"cloud.google.com/go/firestore"
)

func CreateNotification(ctx context.Context, fs *firestore.Client, user models.User, reference models.Reference, message string) error {

	notifColl := fs.Collection("notifications")
	notifRef := notifColl.NewDoc()
	notification := models.Notification{
		ID:        notifRef.ID,
		User:      user,
		Message:   message,
		Reference: reference,
		Read:      false,
		CreatedAt: time.Now(),
	}
	_, err := notifRef.Set(ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

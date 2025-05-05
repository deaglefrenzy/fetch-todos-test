package groups

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/notifications"
	"learnfirestore/utils"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

func SeedComment(createdBy models.User) models.Comments {

	return models.Comments{
		UUID:      uuid.NewString(),
		Text:      utils.GenerateString(10),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}

func CommentsFactory(ctx context.Context, fs *firestore.Client, createdBy models.User, group models.Group, count int) error {

	groupRef := fs.Collection("groups").Doc(group.ID)
	groupDoc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	var members []models.User
	if err := groupDoc.DataTo(&members); err != nil {
		return fmt.Errorf("failed to unmarshal members: %w", err)
	}

	var comments []models.Comments
	if err := groupDoc.DataTo(&comments); err != nil {
		comments = []models.Comments{}
	}

	for range count {
		c := SeedComment(createdBy)
		comments = append(comments, c)
		fmt.Printf("Comment %s created\n", c.Text)
	}
	_, err = groupRef.Update(ctx, []firestore.Update{
		{
			Path:  "comments",
			Value: comments,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update group comments: %w", err)
	}

	for _, m := range members {
		for _, c := range comments {
			reference := models.Reference{
				Type: "comment",
				Data: map[string]interface{}{
					"id":   c.UUID,
					"name": c.Text,
				},
			}
			err = notifications.CreateNotification(ctx, fs, m, reference, "There's a new comment in the group.")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

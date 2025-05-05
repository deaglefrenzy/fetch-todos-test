package groups

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/notifications"
	"learnfirestore/utils"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

func SeedTask(createdBy models.User) models.Tasks {

	return models.Tasks{
		UUID:        uuid.NewString(),
		Title:       utils.GenerateString(3),
		Description: utils.GenerateString(6),
		Priority:    rand.Intn(2) == 1,
		Done:        false,
		DueDate:     time.Now().AddDate(0, 0, 1),
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}
}

func TaskFactory(ctx context.Context, fs *firestore.Client, createdBy models.User, group models.Group, count int) error {

	groupRef := fs.Collection("groups").Doc(group.ID)
	groupDoc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	var members []models.User
	if err := groupDoc.DataTo(&members); err != nil {
		return fmt.Errorf("failed to unmarshal members: %w", err)
	}

	var tasks []models.Tasks
	if err := groupDoc.DataTo(&tasks); err != nil {
		tasks = []models.Tasks{}
	}

	for range count {
		t := SeedTask(createdBy)
		tasks = append(tasks, t)
		fmt.Printf("Task %s created\n", t.Title)
	}
	_, err = groupRef.Update(ctx, []firestore.Update{
		{
			Path:  "tasks",
			Value: tasks,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update group tasks: %w", err)
	}

	for _, m := range members {
		for _, c := range tasks {
			reference := models.Reference{
				Type: "tasks",
				Data: map[string]interface{}{
					"id":   c.UUID,
					"name": c.Title,
				},
			}
			err = notifications.CreateNotification(ctx, fs, m, reference, "There's a new task in the group.")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

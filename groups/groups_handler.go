package groups

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/notifications"
	"learnfirestore/utils"
	"time"

	"cloud.google.com/go/firestore"
)

func SeedGroup() models.Group {
	return models.Group{
		Name:        utils.GenerateString(2),
		Description: utils.GenerateString(4),
		Members:     []models.User{},
		Tasks:       []models.Tasks{},
		Comments:    []models.Comments{},
		CreatedAt:   time.Now(),
	}
}

func GroupFactory(ctx context.Context, fs *firestore.Client, count int, member models.User) ([]models.Group, error) {

	var group []models.Group
	for range count {
		g := SeedGroup()
		groupColl := fs.Collection("groups").NewDoc()
		g.ID = groupColl.ID
		g.Members = append(g.Members, member)
		_, err := groupColl.Set(ctx, g)
		if err != nil {
			return nil, fmt.Errorf("failed to create group: %w", err)
		}
		group = append(group, g)
		fmt.Printf("Group %s (%s) created\n", g.Name, g.ID)

		reference := models.Reference{
			Type: "group",
			Data: map[string]interface{}{
				"id":   g.ID,
				"name": g.Name,
			},
		}
		err = notifications.CreateNotification(ctx, fs, member, reference, "You are added into a group!")
		if err != nil {
			return nil, err
		}
	}
	return group, nil
}

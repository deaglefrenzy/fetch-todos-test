package main

import (
	"learnfirestore/database"
	"learnfirestore/repository"
	"log"
)

func main() {

	ctx, fs, err := database.ConnectToFirestore()
	if err != nil {
		panic(err)
	}

	// newUser, err := users.UserFactory(ctx, fs, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// member := models.User{ID: newUser[0].ID, Name: newUser[0].Name}

	// newGroup, err := groups.GroupFactory(ctx, fs, 1, member)
	// if err != nil {
	// 	panic(err)
	// }

	// err = groups.CommentsFactory(ctx, fs, member, newGroup[0], 1)
	// if err != nil {
	// 	panic(err)
	// }

	// err = groups.TaskFactory(ctx, fs, member, newGroup[0], 1)
	// if err != nil {
	// 	panic(err)
	// }

	go func() {
		if err := repository.WatchGroups(ctx, fs); err != nil {
			log.Printf("Watch error: %v", err)
		}
	}()

	select {}
}

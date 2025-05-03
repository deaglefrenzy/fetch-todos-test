package main

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/google/uuid"

	"google.golang.org/api/option"
)

func generateString(numWords int) string {
	if numWords <= 0 {
		return ""
	}

	loremIpsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	words := strings.Fields(loremIpsum)
	if len(words) == 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	generatedWords := make([]string, numWords)

	for i := 0; i < numWords; i++ {
		randomIndex := r.Intn(len(words))
		generatedWords[i] = words[randomIndex]
	}

	return strings.Join(generatedWords, " ")
}

func seedUser() models.UserCollection {
	return models.UserCollection{
		Name:      generateString(2),
		Email:     generateString(1) + "@" + generateString(1) + ".com",
		CreatedAt: time.Now(),
	}
}

func seedTasks(createdBy models.User, count int) []models.Tasks {
	tasks := make([]models.Tasks, count)

	for i := 0; i < count; i++ {
		tasks[i] = models.Tasks{
			UUID:        uuid.NewString(),
			Title:       generateString(3),
			Description: generateString(6),
			Priority:    rand.Intn(2) == 1,
			Done:        false,
			DueDate:     time.Now().AddDate(0, 0, 1),
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
		}
	}
	return tasks
}

func seedGroup() models.Group {
	return models.Group{
		Name:        generateString(2),
		Description: generateString(4),
		Members:     []models.User{},
		Tasks:       []models.Tasks{},
		Comments:    []models.Comments{},
		CreatedAt:   time.Now(),
	}
}

func seedComments(createdBy models.User, count int) []models.Comments {
	comments := make([]models.Comments, count)

	for i := 0; i < count; i++ {
		comments[i] = models.Comments{
			UUID:      uuid.NewString(),
			Text:      generateString(10),
			CreatedBy: createdBy,
			CreatedAt: time.Now(),
		}
	}
	return comments
}

func seedNotifications(ctx context.Context, fs *firestore.Client, user models.User, count int) {
	notification := make([]models.Notification, count)

	notifColl := fs.Collection("notifications")

	for i := 0; i < count; i++ {
		notifDoc := notifColl.NewDoc()
		notification[i] = models.Notification{
			ID:        notifDoc.ID,
			User:      user,
			Type:      "Common",
			Message:   generateString(7),
			Read:      false,
			CreatedAt: time.Now(),
		}
		_, err := notifDoc.Set(ctx, notification)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	ctx := context.Background() // digunakan untuk action yg memerlukan time

	opt := option.WithCredentialsFile("service.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic(err)
	}

	fs, err := app.Firestore(ctx) // requires connection to firestore database/server
	if err != nil {
		panic(err)
	}

	user := seedUser()
	userColl := fs.Collection("users").NewDoc()
	user.ID = userColl.ID
	_, err = userColl.Set(ctx, user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User %s : %s created\n", user.ID, user.Name)

	member := models.User{ID: user.ID, Name: user.Name}

	groupColl := fs.Collection("groups")
	seedNotifications(ctx, fs, member, 3)

	for i := 0; i < 2; i++ {

		group := seedGroup()

		newDoc := groupColl.NewDoc()
		group.ID = newDoc.ID
		group.Members = append(group.Members, member)
		group.Comments = seedComments(member, 3)
		group.Tasks = seedTasks(member, 10)

		_, err := newDoc.Set(ctx, group)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Group %d %s created with %d tasks\n", i+1, group.Name, len(group.Tasks))
	}
}

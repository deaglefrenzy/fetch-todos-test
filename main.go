package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type Group struct {
	ID          string     `firestore:"-"`
	Name        string     `firestore:"name" json:"name"`
	Description string     `firestore:"description" json:"description"`
	Members     []User     `firestore:"members" json:"members"`
	Tasks       []Tasks    `firestore:"tasks" json:"tasks,omitempty"`
	Comments    []Comments `firestore:"comments" json:"comments,omitempty"`
	CreatedAt   time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt   *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type User struct {
	ID           string         `firestore:"-"` //omitempty versi firestore?
	Name         string         `firestore:"name" json:"name"`
	Email        string         `firestore:"email" json:"email"`
	Notification []Notification `firestore:"notification" json:"notification,omitempty"`
	CreatedAt    time.Time      `firestore:"createdAt" json:"created_at"`
	DeletedAt    *time.Time     `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type Tasks struct {
	ID          string     `firestore:"-"`
	Title       string     `firestore:"title" json:"title"`
	Description string     `firestore:"description" json:"description"`
	Priority    bool       `firestore:"priority" json:"priority"`
	Done        bool       `firestore:"done" json:"done"`
	DueDate     time.Time  `firestore:"dueDate" json:"due_date"`
	CreatedBy   User       `firestore:"createdBy" json:"created_by"`
	CreatedAt   time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt   *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type Comments struct {
	Text      string     `firestore:"text" json:"text"`
	CreatedBy User       `firestore:"createdBy" json:"created_by"`
	CreatedAt time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type Notification struct {
	Type        string    `firestore:"type" json:"type"`
	Message     string    `firestore:"message" json:"message"`
	ReferenceID string    `firestore:"referenceID" json:"reference_id"`
	Read        bool      `firestore:"read" json:"read"`
	CreatedAt   time.Time `firestore:"createdAt" json:"created_at"`
}

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

func createTask(createdBy User) Tasks {
	return Tasks{
		Title:       generateString(3),
		Description: generateString(6),
		Priority:    rand.Intn(2) == 1,
		Done:        false,
		DueDate:     time.Now().AddDate(0, 0, 1),
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}
}

func createGroup() Group {
	return Group{
		Name:        generateString(2),
		Description: generateString(4),
		Members:     []User{},
		Tasks:       []Tasks{},
		Comments:    []Comments{},
		CreatedAt:   time.Now(),
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

	for i := 0; i < 10; i++ {

		seedUser := User{Name: "Random User", Email: "random-user@email.com", CreatedAt: time.Now()}
		group := createGroup()

		groupColl := fs.Collection("groups").NewDoc()
		group.ID = groupColl.ID

		for j := 0; j < 10; j++ {
			task := createTask(seedUser)

			taskColl := groupColl.Collection("tasks").NewDoc()
			task.ID = taskColl.ID

			group.Tasks = append(group.Tasks, task)
		}

		_, err := groupColl.Set(ctx, group)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Group %d %s created with %d tasks\n", i+1, group.Name, len(group.Tasks))
	}

}

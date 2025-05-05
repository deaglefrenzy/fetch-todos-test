package learn

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type User struct {
	ID           string     `firestore:"-"` //omitempty versi firestore?
	Name         string     `firestore:"name" json:"name"`
	Email        string     `firestore:"email" json:"email"`
	Notification string     `firestore:"notification" json:"notification,omitempty"`
	CreatedAt    time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt    *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

func learn() {
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

	fmt.Println(fs)

	coll := fs.Collection("users")
	doc := coll.Doc("Fv9jTTr9l65BgpvGdVBe")

	snap, err := doc.Get(ctx)
	if err != nil {
		panic(err)
	}

	var user User
	fmt.Println(user)
	err = snap.DataTo(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)

	doc.Update(ctx, []firestore.Update{
		{Path: "deletedAt", Value: time.Now()},
	})

	user.DeletedAt = nil
	doc.Set(ctx, user, firestore.MergeAll)

	// data := snap.Data()
	// id := doc.ID
	// name := (data["name"]).(string) //this is typecasting (forced, panic if fails)

	// newdata := User{ID: id, Name: name}

	// fmt.Println(newdata)

	//fmt.Println(name)

	// filter, _ := coll.Where("name", "==", "user1").Documents(ctx).GetAll()
	// if len(filter) == 0 {
	// 	panic("no results")
	// }

	// for _, val := range filter {
	// 	fmt.Println(val.Data()["name"])
	// }

	// iter := filter.Documents(ctx) // requires connection
	// snap, _ := iter.GetAll()
	// for _, val := range snap {
	// 	fmt.Println(val.Data()["name"])
	// }
	// c := char{ID: 1, Name: "nama"}
	// c.toString()
}

// func toString() string { // this function belongs to char struct therefore it has access to all of its attributes
// 	return "My ID" + this.ID + " and my name is" + this.Name
// }

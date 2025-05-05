package learn

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func conc() {
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

	coll := fs.Collection("groups")

	snap := coll.Snapshots(ctx)

	dataChan := make(chan string, 5)

	go func() {
		for {
			fmt.Println("CHECKING CHANGES")
			qs, err := snap.Next()
			if err != nil {
				panic(err)
			}

			for _, v := range qs.Changes {
				fmt.Printf("Document change: %v\n", v.Kind)

				if v.Kind == firestore.DocumentAdded {
					dataChan <- fmt.Sprintf("added:%s", v.Doc.Ref.ID)
				} else if v.Kind == firestore.DocumentModified {
					dataChan <- fmt.Sprintf("modified:%s", v.Doc.Ref.ID)
				} else if v.Kind == firestore.DocumentRemoved {
					dataChan <- fmt.Sprintf("removed:%s", v.Doc.Ref.ID)
				}
			}
		}
	}()

	go func() {
		for {
			fmt.Println("WRITING CHANGES")
			data := <-dataChan
			fmt.Println(data)
		}
	}()

	for {
	}
}

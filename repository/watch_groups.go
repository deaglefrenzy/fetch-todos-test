package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

func WatchGroups(ctx context.Context, fs *firestore.Client) error {
	interval := 5 * time.Second

	coll := fs.Collection("groups")
	snap := coll.Snapshots(ctx)

	fmt.Println("Start watching Groups Collection...")

	for {
		qs, err := snap.Next()
		if err != nil {
			panic(err)
		}

		for _, v := range qs.Changes {
			groupID := v.Doc.Ref.ID

			if v.Kind == firestore.DocumentAdded {
				fmt.Printf("added:%s\n", groupID)
			} else if v.Kind == firestore.DocumentModified {
				fmt.Printf("modified:%s\n", groupID)
			} else if v.Kind == firestore.DocumentRemoved {
				fmt.Printf("removed:%s\n", groupID)
			}
		}

		time.Sleep(interval)
	}
}

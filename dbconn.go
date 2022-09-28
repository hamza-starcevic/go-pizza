package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

//Connect to firestore database

// ! Establish a connection to the firestore database
func Init() (*firestore.Client, context.Context) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "pizza-847ab")
	if err != nil {
		log.Fatalf("Failed to create a firestore client: %v", err)
	}
	return client, ctx
}

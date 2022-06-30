package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

var ctx = context.Background()

func FirebaseApp() (c *firebase.App, err error) {
	opt := option.WithCredentialsFile("credentials.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return app, err
}

func GetFirestore() (c *firestore.Client, err error) {
	app, err := FirebaseApp()
	if err != nil {
		fmt.Errorf("error registering FirebaseApp: %v", err)
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, err
}

func GetFirebaseAuth() (c *auth.Client, err error) {
	app, err := FirebaseApp()
	if err != nil {
		fmt.Errorf("error registering FirebaseApp: %v", err)
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, err
}

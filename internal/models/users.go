package models

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

type User struct {
	ID    string `firestore:"-"`
	Email string `firestore:"email"`
}

type UserModel struct {
	DB   *firestore.Client
	Auth *firebase.App
}

func (m *UserModel) CreateUser(email string, uid string) error {
	ctx := context.Background()
	// add the uid as the doc id in firestore
	_, err := m.DB.Collection("users").Doc(uid).Set(ctx, User{Email: email})
	if err != nil {
		return err
	}
	return nil
}

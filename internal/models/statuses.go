package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Status struct {
	Text  string `firestore:"text"`
	Color string `firestore:"color"`
}

type StatusModel struct {
	DB *firestore.Client
}

func (s *StatusModel) Insert(userId, boardId, groupId, taskId, name, color string) error {
	ctx := context.Background()
	status := Status{
		Text:  name,
		Color: color,
	}
	_, _, err := s.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Collection("tasks").Doc(taskId).Collection("statuses").Add(ctx, status)
	if err != nil {
		return err
	}
	return nil
}

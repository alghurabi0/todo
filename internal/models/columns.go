package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Column struct {
	ID      string `firestore:"-"`
	BoardId string `firestore:"-"`
	Name    string `firestore:"name"`
	Type    string `firestore:"type"`
}

type ColumnModel struct {
	DB *firestore.Client
}

func (m *ColumnModel) Insert(userId, boardId, name, colType string) (string, error) {
	ctx := context.Background()
	col := Column{
		Name:    name,
		Type:    colType,
		BoardId: boardId,
	}
	doc, _, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("columns").Add(ctx, col)
	if err != nil {
		return "", err
	}
	col.ID = doc.ID
	return doc.ID, nil
}

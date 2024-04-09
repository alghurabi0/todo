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

	// update columnOrder array in the board document
	boardDoc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Get(ctx)
	if err != nil {
		return "", err
	}
	var board Board
	if err := boardDoc.DataTo(&board); err != nil {
		return "", err
	}
	board.ColumnOrder = append(board.ColumnOrder, doc.ID)
	_, err = m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Update(ctx, []firestore.Update{
		{Path: "column_order", Value: board.ColumnOrder},
	})
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}
func (m *ColumnModel) GetColumnOrder(userId, boardId string) ([]string, error) {
	ctx := context.Background()
	doc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Get(ctx)
	if err != nil {
		return nil, err
	}
	var board Board
	if err := doc.DataTo(&board); err != nil {
		return nil, err
	}
	return board.ColumnOrder, nil
}
func (m *ColumnModel) Reorder(userId string, boardId string, order []string) error {
	ctx := context.Background()
	_, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Update(ctx, []firestore.Update{
		{Path: "column_order", Value: order},
	})
	if err != nil {
		return err
	}
	return nil
}

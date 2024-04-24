package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type Group struct {
	ID            string    `firestore:"-"`
	BoardId       string    `firestore:"-"`
	Name          string    `firestore:"name"`
	Order         int       `firestore:"order"`
	Tasks         Tasks     `firestore:"-"`
	LastTaskOrder int       `firestore:"last_task_order"`
	CreatedAt     time.Time `firestore:"created_at"`
}

type GroupModel struct {
	DB *firestore.Client
}

func (m *GroupModel) Insert(userId string, boardId string, name string) (*Group, error) {
	ctx := context.Background()
	order, err := m.GetLastGroupOrder(userId, boardId)
	if err != nil {
		return nil, err
	}
	group := Group{
		Name:          name,
		Order:         order + 1,
		LastTaskOrder: 0,
		CreatedAt:     time.Now(),
	}
	doc, _, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Add(ctx, group)
	if err != nil {
		return nil, err
	}
	err = m.UpdateLastGroupOrder(userId, boardId, order+1)
	if err != nil {
		return nil, err
	}
	group.ID = doc.ID
	return &group, nil
}
func (m *GroupModel) GetLastGroupOrder(userId string, boardId string) (int, error) {
	ctx := context.Background()
	doc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Get(ctx)
	if err != nil {
		return 0, err
	}
	var board Board
	if err := doc.DataTo(&board); err != nil {
		return 0, err
	}
	return board.LastGroupOrder, nil
}
func (m *GroupModel) UpdateLastGroupOrder(userId string, boardId string, order int) error {
	ctx := context.Background()
	_, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Update(ctx, []firestore.Update{
		{Path: "last_group_order", Value: order},
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *GroupModel) UpdateName(userId, boardId, groupId, groupName string) error {
	ctx := context.Background()
	_, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Update(ctx, []firestore.Update{
		{Path: "name", Value: groupName},
	})
	if err != nil {
		return err
	}
	return nil
}

type Groups []Group

func (g Groups) Len() int {
	return len(g)
}
func (g Groups) Less(i, j int) bool {
	return g[i].Order < g[j].Order
}
func (g Groups) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

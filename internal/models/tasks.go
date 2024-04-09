package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type Task struct {
	ID           string                 `firestore:"-"`
	BoardId      string                 `firestore:"-"`
	GroupId      string                 `firestore:"-"`
	Content      string                 `firestore:"content"`
	Order        int                    `firestore:"order"`
	CreatedAt    time.Time              `firestore:"created_at"`
	ColumnValues map[string]interface{} `firestore:"column_values"`
	ColumnOrder  []string               `firestore:"-"`
}

type TaskModel struct {
	DB *firestore.Client
}

func (m *TaskModel) Insert(userId string, boardId string, groupId string, content string) (*Task, error) {
	ctx := context.Background()
	order, err := m.GetLastTaskOrder(userId, boardId, groupId)
	if err != nil {
		return nil, err
	}
	task := Task{
		Content:   content,
		CreatedAt: time.Now(),
		Order:     order + 1,
	}
	doc, _, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Collection("tasks").Add(ctx, task)
	if err != nil {
		return nil, err
	}
	err = m.UpdateLastTaskOrder(userId, boardId, groupId, order+1)
	if err != nil {
		return nil, err
	}
	task.ID = doc.ID
	return &task, nil
}
func (m *TaskModel) Get(id string) (*Task, error) {
	ctx := context.Background()
	doc, err := m.DB.Collection("tasks").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := doc.DataTo(&task); err != nil {
		return nil, err
	}
	task.ID = doc.Ref.ID
	return &task, nil
}
func (m *TaskModel) Delete(userId string, boardId string, groupId string, taskId string) error {
	ctx := context.Background()
	_, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Collection("tasks").Doc(taskId).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (m *TaskModel) GetLastTaskOrder(userId string, boardId string, groupId string) (int, error) {
	ctx := context.Background()
	doc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Get(ctx)
	if err != nil {
		return 0, err
	}
	var group Group
	if err := doc.DataTo(&group); err != nil {
		return 0, err
	}
	return group.LastTaskOrder, nil
}
func (m *TaskModel) UpdateLastTaskOrder(userId string, boardId string, groupId string, newOrder int) error {
	ctx := context.Background()
	_, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Update(ctx, []firestore.Update{
		{Path: "last_task_order", Value: newOrder},
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *TaskModel) Swap(userId string, boardId string, groupId string,
	swappedId string, swappedOrder int, targetId string, targetOrder int) error {
	ctx := context.Background()
	err := m.DB.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// update the swapped task without getting the document, since we already have the order
		err := tx.Update(m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Collection("tasks").Doc(swappedId), []firestore.Update{
			{Path: "order", Value: targetOrder},
		})
		if err != nil {
			return err
		}
		// update the target task without getting the document, since we already have the order
		err = tx.Update(m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupId).Collection("tasks").Doc(targetId), []firestore.Update{
			{Path: "order", Value: swappedOrder},
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

type Tasks []Task

func (t Tasks) Len() int {
	return len(t)
}
func (t Tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Tasks) Less(i, j int) bool {
	return t[i].Order < t[j].Order
}

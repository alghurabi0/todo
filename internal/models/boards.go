package models

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Board struct {
	ID             string            `firestore:"-"`
	Title          string            `firestore:"title"`
	Description    string            `firestore:"description"`
	CreatedAt      time.Time         `firestore:"created_at"`
	Groups         Groups            `firestore:"-"`
	LastGroupOrder int               `firestore:"last_group_order"`
	Columns        map[string]Column `firestore:"-"`
	ColumnOrder    []string          `firestore:"column_order"`
}

type BoardModel struct {
	DB *firestore.Client
}

func (m *BoardModel) Insert(title string, description string, userId string) (string, error) {
	ctx := context.Background()
	ref, _, err := m.DB.Collection("users").Doc(userId).Collection("boards").Add(ctx, Board{
		Title:          title,
		Description:    description,
		LastGroupOrder: 0,
		ColumnOrder:    []string{},
		CreatedAt:      time.Now(),
	})
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (m *BoardModel) Get(userId string, boardId string) (*Board, error) {
	ctx := context.Background()
	doc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Get(ctx)
	if err != nil {
		return nil, err
	}
	var board Board
	if err := doc.DataTo(&board); err != nil {
		return nil, err
	}
	board.Columns = make(map[string]Column)
	// Get all columns for the board
	for _, colId := range board.ColumnOrder {
		colDoc, err := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("columns").Doc(colId).Get(ctx)
		if err != nil {
			return nil, err
		}
		var col Column
		if err := colDoc.DataTo(&col); err != nil {
			return nil, err
		}
		col.ID = colDoc.Ref.ID
		col.BoardId = boardId
		board.Columns[col.ID] = col
	}
	board.ID = doc.Ref.ID
	// Get all groups for the board
	groupIter := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Documents(ctx)
	for {
		groupDoc, err := groupIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var group Group
		if err := groupDoc.DataTo(&group); err != nil {
			return nil, err
		}
		group.ID = groupDoc.Ref.ID
		// Get all tasks for the group
		taskIter := m.DB.Collection("users").Doc(userId).Collection("boards").Doc(boardId).Collection("groups").Doc(groupDoc.Ref.ID).Collection("tasks").Documents(ctx)
		for {
			taskDoc, err := taskIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			var task Task
			if err := taskDoc.DataTo(&task); err != nil {
				return nil, err
			}
			task.ID = taskDoc.Ref.ID
			task.BoardId = boardId
			task.GroupId = groupDoc.Ref.ID
			group.Tasks = append(group.Tasks, task)
		}
		group.BoardId = boardId
		// Sort the tasks by order
		sort.Sort(group.Tasks)
		board.Groups = append(board.Groups, group)
	}
	// Sort the groups by order
	sort.Sort(board.Groups)
	return &board, nil
}

func (m *BoardModel) GetAll(userId string) (*[]Board, error) {
	ctx := context.Background()
	boardIter := m.DB.Collection("users").Doc(userId).Collection("boards").Documents(ctx)
	var boards []Board
	for {
		doc, err := boardIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var board Board
		if err := doc.DataTo(&board); err != nil {
			return nil, err
		}
		board.ID = doc.Ref.ID
		boards = append(boards, board)
	}
	return &boards, nil
}

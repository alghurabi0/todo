// each task has a messages (updates) scroll window
// each message has a user tag
// messages can have replies
// reply to a reply is another reply with a tag to the user who first replied
package models

import "time"

type Message struct {
	ID        string    `firestore:"-"`
	TaskId    string    `firestore:"-"`
	Content   string    `firestore:"content"`
	UserId    string    `firestore:"user_id"`
	CreatedAt time.Time `firestore:"created_at"`
	Replies   []Message `firestore:"-"`
}

// start by creating a Messages collection in the task document

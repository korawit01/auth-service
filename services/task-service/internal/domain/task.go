package domain

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	RowId       string     `db:"row_id" json:"row_id"`
	UserID      string     `db:"user_id" json:"-"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description,omitempty"`
	Status      TaskStatus `db:"status" json:"status"`
	DueDate     *time.Time `db:"due_date" json:"dueDate,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}

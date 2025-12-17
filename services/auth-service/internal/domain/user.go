package domain

import "time"

type User struct {
	RowId     string    `db:"row_id" json:"row_id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

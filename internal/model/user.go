package model

import (
	"database/sql"
	"time"
)

type User struct {
	UserId    int            `json:"id"`
	Username  string         `json:"name"`
	Email     sql.NullString `json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
}

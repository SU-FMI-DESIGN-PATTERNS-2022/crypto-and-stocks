package user_repository

import "database/sql"

type User struct {
	ID        int64
	UserID    sql.NullInt64
	Name      sql.NullString
	IsBot     bool
	CreatorID sql.NullInt64
	Amount    float64
}

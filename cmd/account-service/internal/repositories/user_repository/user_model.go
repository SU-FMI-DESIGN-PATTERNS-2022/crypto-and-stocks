package user_repository

import "database/sql"

type User struct {
	ID        int64          `db:"id"`
	UserID    sql.NullInt64  `db:"user_id"`
	Name      sql.NullString `db:"name"`
	IsBot     bool           `db:"is_bot"`
	CreatorID sql.NullInt64  `db:"creator_id"`
	Amount    float64        `db:"amount"`
}

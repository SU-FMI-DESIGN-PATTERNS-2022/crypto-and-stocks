package user_repository

import (
	"database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type UserTable struct {
	instance *sql.DB
}

func NewUserTable(db *sql.DB) *UserTable {
	return &UserTable{
		instance: db,
	}
}

func (db *UserTable) MergeUserOrders(id int64) error {
	rows, err := db.instance.Query(selectAllWhereCreatorIdSQL, id)

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			pq.Array(&user.Orders),
			&user.IsBot,
			&user.CreatorID,
			&user.Amount,
		)

		if err != nil {

			return err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return err
	}

	var orders []int64
	for _, u := range users {
		orders = append(orders, u.Orders...)
		_, deleteErr := db.instance.Exec(deleteUserWhereIdSQL, u.ID)

		if deleteErr != nil {
			return deleteErr
		}
	}

	row := db.instance.QueryRow(selectOrdersWhereIdSQL, id)
	var userOrders []int64
	ordersErr := row.Scan(pq.Array(&userOrders))

	if ordersErr != nil || row.Err() != nil {
		return ordersErr
	}

	orders = append(orders, userOrders...)

	_, updateErr := db.instance.Exec(updateUserOrdersWhereIdSQL, pq.Array(orders), id)

	return updateErr
}

// might need to calculate amount of bots as well
func (db *UserTable) GetCurrentAmount(id int64) (float64, error) {
	row := db.instance.QueryRow(selectAmountWhereIdSQL, id)

	var amount float64
	err := row.Scan(&amount)

	if err != nil || row.Err() != nil {
		return 0, err
	}

	return amount, nil
}

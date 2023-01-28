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

func (db *UserTable) CreateUser(userId int64, name string) error {
	_, err := db.instance.Exec(createUserSQL,
		userId,
		name,
		pq.Array(make([]int64, 0)),
		false,
		nil,
		0,
	)

	return err
}

func (db *UserTable) CreateBot(creatorID int64, amount float64) error {
	_, err := db.instance.Exec(createBotSQL,
		nil,
		nil,
		pq.Array(make([]int64, 0)),
		true,
		nil,
		amount,
	)

	return err
}

func (db *UserTable) AddOrder(userId int64, orderId int64) error {
	row := db.instance.QueryRow(selectOrdersWhereIdSQL, userId)
	var orders []int64
	ordersErr := row.Scan(pq.Array(&orders))

	if ordersErr != nil || row.Err() != nil {
		return ordersErr
	}

	orders = append(orders, orderId)
	_, updateErr := db.instance.Exec(updateUserOrdersWhereIdSQL, pq.Array(orders), userId)

	return updateErr
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

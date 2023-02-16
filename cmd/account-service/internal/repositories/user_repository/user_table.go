package user_repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserTable struct {
	instance *sqlx.DB
}

func NewUserTable(db *sqlx.DB) *UserTable {
	return &UserTable{
		instance: db,
	}
}

func (table *UserTable) CreateUser(userId int64, name string) error {
	_, err := table.instance.Exec(createUserSQL,
		userId,
		name,
		false,
		nil,
		0,
	)

	return err
}

func (table *UserTable) GetUserById(id int64) (User, error) {
	var user User
	err := table.instance.Get(&user, selectUserWhereIdSQL, id)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (table *UserTable) GetUserByUserId(id int64) (User, error) {
	var user User
	err := table.instance.Get(&user, selectUserWhereUserIdSQL, id)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (table *UserTable) GetUserAmount(id int64) (float64, error) {
	var amount float64
	err := table.instance.Get(&amount, selectUserAmountWhereIdSQL, id)

	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (table *UserTable) CreateBot(creatorID int64, amount float64) error {
	_, err := table.instance.Exec(createBotSQL,
		nil,
		nil,
		true,
		creatorID,
		amount,
	)

	return err
}

func (table *UserTable) UpdateUserAmount(id int64, amount float64) error {
	_, err := table.instance.Exec(updateUserAmountSQL, amount, id)
	return err
}

func (table *UserTable) DeleteUserById(id int64) error {
	_, deleteErr := table.instance.Exec(deleteUserWhereIdSQL, id)
	return deleteErr
}

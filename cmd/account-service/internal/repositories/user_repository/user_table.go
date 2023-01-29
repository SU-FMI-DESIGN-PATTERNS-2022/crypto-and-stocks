package user_repository

import (
	"database/sql"

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
		false,
		nil,
		0,
	)

	return err
}

func (db *UserTable) CreateBot(creatorID int64, amount float64) error {
	//TODO: check if the amount is valid and if it is subtract it from the user amount
	_, err := db.instance.Exec(createBotSQL,
		nil,
		nil,
		true,
		creatorID,
		amount,
	)

	return err
}

func (db *UserTable) MergeUserAndBot(id int64) error {
	//TODO: after merge update user amount by bot amount
	row := db.instance.QueryRow(selectUserWhereIdSQL, id)

	var bot User
	userErr := row.Scan(
		&bot.ID,
		&bot.Name,
		&bot.UserID,
		&bot.IsBot,
		&bot.CreatorID,
		&bot.Amount,
	)

	if userErr != nil {
		return userErr
	}

	_, ordersErr := db.instance.Exec(updateOrdersAfterMergeSQL, bot.CreatorID, id)

	if ordersErr != nil {
		return ordersErr
	}

	_, deleteErr := db.instance.Exec(deleteUserWhereIdSQL, id)

	return deleteErr
}

func (db *UserTable) MergeAllUserOrders(id int64) error {
	botsRows, err := db.instance.Query(selectAllWhereCreatorIdSQL, id)

	defer botsRows.Close()

	var bots []int64
	for botsRows.Next() {
		var botId int64
		err := botsRows.Scan(&botId)

		if err != nil {
			return err
		}

		bots = append(bots, botId)
	}

	if botsRows.Err() != nil {
		return err
	}

	for _, b := range bots {
		//TODO: after merge update user amount by bot amount
		_, ordersErr := db.instance.Exec(updateOrdersAfterMergeSQL, id, b)

		if ordersErr != nil {
			return ordersErr
		}

		_, deleteErr := db.instance.Exec(deleteUserWhereIdSQL, b)

		if deleteErr != nil {
			return deleteErr
		}
	}

	return err
}

func (db *UserTable) GetAmountByUserId(id int64) (float64, error) {
	//TODO
	return 0, nil
}

package user_repository

import (
	"errors"
	"fmt"
	"math"

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

func (table *UserTable) CreateBot(creatorID int64, amount float64) error {
	var creator User
	userErr := table.instance.Get(&creator, selectUserWhereIdSQL, creatorID)

	if userErr != nil {
		return userErr
	}

	if creator.Amount < amount {
		return errors.New("Not enough amount to create bot")
	}

	if creator.IsBot {
		return errors.New("Bots can't have their own bots")
	}

	_, updateErr := table.instance.Exec(updateUserAmountSQL, math.Round((creator.Amount-amount)*100)/100, creator.ID)

	if updateErr != nil {
		return updateErr
	}

	_, err := table.instance.Exec(createBotSQL,
		nil,
		nil,
		true,
		creatorID,
		amount,
	)

	return err
}

func (table *UserTable) MergeUserAndBot(id int64) error {
	var bot User
	botErr := table.instance.Get(&bot, selectUserWhereIdSQL, id)

	if botErr != nil {
		return botErr
	}

	if !bot.IsBot {
		return errors.New("Can't merge 2 users - one must be bot")
	}

	var user User
	userErr := table.instance.Get(&user, selectUserWhereIdSQL, bot.CreatorID)

	if userErr != nil {
		return userErr
	}

	_, ordersErr := table.instance.Exec(updateOrdersAfterMergeSQL, bot.CreatorID, id)

	if ordersErr != nil {
		return ordersErr
	}

	_, amountErr := table.instance.Exec(updateUserAmountSQL, math.Round((user.Amount+bot.Amount)*100)/100, user.ID)

	if amountErr != nil {
		return amountErr
	}

	_, deleteErr := table.instance.Exec(deleteUserWhereIdSQL, id)

	return deleteErr
}

func (table *UserTable) MergeAllUserOrders(id int64) error {
	var user User
	userErr := table.instance.Get(&user, selectUserWhereIdSQL, id)
	if userErr != nil {
		return userErr
	}

	if user.IsBot {
		return errors.New("Bot can't merge with other users")
	}

	var bots []User
	err := table.instance.Select(&bots, selectAllWhereCreatorIdSQL, id)

	if err != nil {
		return err
	}

	for _, b := range bots {
		//TODO: fix bug where some queries are not executed
		fmt.Println(user.Amount, b.Amount)
		_, amountErr := table.instance.Exec(updateUserAmountSQL, math.Round((user.Amount+b.Amount)*100)/100, user.ID)

		if amountErr != nil {
			return amountErr
		}

		_, ordersErr := table.instance.Exec(updateOrdersAfterMergeSQL, id, b.ID)

		if ordersErr != nil {
			return ordersErr
		}

		_, deleteErr := table.instance.Exec(deleteUserWhereIdSQL, b.ID)

		if deleteErr != nil {
			return deleteErr
		}
	}

	return err
}

func (table *UserTable) GetAmountByUserId(id int64) (float64, error) {
	//TODO
	return 0, nil
}

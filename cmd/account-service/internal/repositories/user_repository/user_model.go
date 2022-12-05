package user_repository

type User struct {
	ID        int64
	Name      string
	Orders    []int64
	IsBot     bool
	CreatorID int64
	Amount    float64
}

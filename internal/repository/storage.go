package repository

type Storage interface {
	UserRepository
	OrderRepository
	WithdrawalRepository
}

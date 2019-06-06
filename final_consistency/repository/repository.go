package repository

type TransactionMsg interface {
}

type Repository interface {
	TransactionMsg
}

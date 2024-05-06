package entity

import "time"

type TransactionID int64

type Transaction struct {
	ID               TransactionID `json:"id" db:"id"`
	ReceivingUserID  UserID        `json:"receivingUserId" db:"receiving_user_id"`
	SendingUserID    UserID        `json:"sendingUserId" db:"sending_user_id"`
	TransactionPoint int           `json:"transactionPoint" db:"transaction_point"`
	TransactionAt    time.Time     `json:"transactionAt" db:"transaction_at"`
}

type Transactions []*Transaction

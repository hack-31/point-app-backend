package entity

import (
	"time"

	"github.com/hack-31/point-app-backend/domain/model"
)

type TransactionID int64

type Transaction struct {
	ID               TransactionID `json:"id" db:"id"`
	ReceivingUserID  model.UserID  `json:"receivingUserId" db:"receiving_user_id"`
	SendingUserID    model.UserID  `json:"sendingUserId" db:"sending_user_id"`
	TransactionPoint int           `json:"transactionPoint" db:"transaction_point"`
	TransactionAt    time.Time     `json:"transactionAt" db:"transaction_at"`
}

type Transactions []*Transaction

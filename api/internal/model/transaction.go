package model

import (
	"time"
)

type Transaction struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	TargetAccount *Account  `json:"targetAccount"`
	SourceAccount *Account  `json:"sourceAccount"`
	ReferenceID   string    `json:"referenceID"`
	Amount        Money     `json:"amount"`
	Description   string    `json:"description"`
}

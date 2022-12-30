package pkg

import "github.com/arikkfir/greenstar/common"

type Transaction struct {
	Date            string          `json:"date"`
	TargetAccountID string          `json:"targetAccountID"`
	SourceAccountID string          `json:"sourceAccountID"`
	ReferenceID     string          `json:"referenceID"`
	Amount          common.Currency `json:"amount"`
	Description     string          `json:"description"`
}

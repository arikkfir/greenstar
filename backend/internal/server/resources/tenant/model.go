// Code generated by greenstar scripts; DO NOT EDIT.

package tenant

import (
	"time"
)

type Tenant struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DisplayName string    `json:"displayName"`
}

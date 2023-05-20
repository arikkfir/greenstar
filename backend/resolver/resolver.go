package resolver

import (
	"github.com/arik-kfir/greenstar/backend/services"
)

type Resolver struct {
	TenantsService      *services.TenantsService
	AccountsService     *services.AccountsService
	TransactionsService *services.TransactionsService
	OperationsService   *services.OperationsService
}

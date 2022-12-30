package internal

import (
	"fmt"
	"github.com/arikkfir/greenstar/common"
	"github.com/arikkfir/greenstar/common/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	// language=Cypher
	deleteAccountCypher = `// delete account
MATCH (account:Account {accountID: $accountID})
DETACH DELETE account`
)

func (s *Server) deleteAccount(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if id == "" {
		c.Negotiate(http.StatusBadRequest, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: fmt.Errorf("ID is required")})
		return
	}

	session := common.GetNeo4jSession(ctx, s.Neo4j)
	defer session.Close(ctx)

	deleteAccountParams := map[string]interface{}{
		"accountID": id,
	}
	if _, err := session.Run(ctx, deleteAccountCypher, deleteAccountParams); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete account")
		c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: gin.H{}})
}

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
	updateAccountCypher = `// update account
MATCH (account:Account {accountID: $accountID})
SET account.displayName = $displayName`

	// language=Cypher
	relinkAccountCypher = `// relink account
MATCH (child:Account {accountID:$childAccountID})-[r:ChildOf]->(oldParent:Account)
WHERE oldParent.accountID <> $parentAccountID
DELETE r
WITH child
MATCH (newParent:Account {accountID:$parentAccountID})
CREATE (child)-[:ChildOf]->(newParent)`

	// language=Cypher
	unlinkAccountCypher = `// relink account
MATCH (child:Account {accountID:$childAccountID})-[r:ChildOf]->(oldParent:Account)
DELETE r`
)

func (s *Server) putAccount(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if id == "" {
		c.Negotiate(http.StatusBadRequest, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: fmt.Errorf("ID is required")})
		return
	}

	input := Account{}
	if err := c.ShouldBind(&input); err != nil {
		ginutil.RenderError(c, http.StatusBadRequest, err)
		return
	}
	account := AccountWithID{ID: id, Account: input}

	session := common.GetNeo4jSession(ctx, s.Neo4j)
	defer session.Close(ctx)

	updateAccountParams := map[string]interface{}{
		"accountID":   account.ID,
		"displayName": account.DisplayName,
	}
	if _, err := session.Run(ctx, updateAccountCypher, updateAccountParams); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update account")
		c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
		return
	}

	if account.ParentID != "" {
		relinkAccountParams := map[string]interface{}{
			"childAccountID":  account.ID,
			"parentAccountID": account.ParentID,
		}
		if _, err := session.Run(ctx, relinkAccountCypher, relinkAccountParams); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to relink account")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
			return
		}
	} else {
		unlinkAccountParams := map[string]interface{}{
			"childAccountID": account.ID,
		}
		if _, err := session.Run(ctx, unlinkAccountCypher, unlinkAccountParams); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to unlink account")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
			return
		}
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: account})
}

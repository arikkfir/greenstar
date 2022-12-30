package internal

import (
	"github.com/arikkfir/greenstar/common"
	"github.com/arikkfir/greenstar/common/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	// language=Cypher
	createAccountCypher = `// create account
CREATE (account:Account {accountID: $accountID, displayName: $displayName})`

	// language=Cypher
	linkToParentAccountCypher = `// link to parent account
MATCH (child:Account {id: $childAccountID}), (parent:Account {id: $parentAccountID})
CREATE (child)-[:ChildOf]->(parent)`
)

func (s *Server) createAccount(c *gin.Context) {
	ctx := c.Request.Context()

	input := Account{}
	if err := c.ShouldBind(&input); err != nil {
		ginutil.RenderError(c, http.StatusBadRequest, err)
		return
	}
	account := AccountWithID{ID: uuid.NewString(), Account: input}

	session := common.GetNeo4jSession(ctx, s.Neo4j)
	defer session.Close(ctx)

	createAccountParams := map[string]interface{}{
		"accountID":   account.ID,
		"displayName": account.DisplayName,
	}
	if _, err := session.Run(ctx, createAccountCypher, createAccountParams); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create account")
		c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
		return
	}

	if account.ParentID != "" {
		linkAccountParams := map[string]interface{}{
			"childAccountID":  account.ID,
			"parentAccountID": account.ParentID,
		}
		if _, err := session.Run(ctx, linkToParentAccountCypher, linkAccountParams); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to link accounts")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
			return
		}
	}

	c.Negotiate(http.StatusCreated, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: account})
}

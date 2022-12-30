package internal

import (
	"github.com/arikkfir/greenstar/common"
	"github.com/arikkfir/greenstar/common/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	// language=Cypher
	getAccountsCypher = `// get accounts
MATCH (account:Account)
OPTIONAL MATCH (account)-[:ChildOf]->(parent:Account)
RETURN account, parent`
)

func (s *Server) getAccounts(c *gin.Context) {
	ctx := c.Request.Context()

	session := common.GetNeo4jSession(ctx, s.Neo4j)
	defer session.Close(ctx)

	result, err := session.Run(ctx, getAccountsCypher, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get accounts")
		c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
		return
	}

	var accounts []AccountWithID
	for result.Next(ctx) {
		record := result.Record()
		accountVal, ok := record.Get("account")
		if !ok {
			log.Ctx(ctx).Error().Err(err).Msg("Missing required 'account' value in record")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
			return
		}

		parentVal, ok := record.Get("parent")
		if !ok {
			log.Ctx(ctx).Error().Err(err).Msg("Missing required 'parent' value in record")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes})
			return
		}

		accountNode := accountVal.(neo4j.Node)
		var parentID string
		if parentNode, ok := parentVal.(neo4j.Node); ok {
			parentID = parentNode.Props["accountID"].(string)
		}

		account := AccountWithID{
			ID: accountNode.Props["accountID"].(string),
			Account: Account{
				DisplayName: accountNode.Props["displayName"].(string),
				ParentID:    parentID,
				Labels:      nil, // TODO: fetch labels
				Annotations: nil, // TODO: fetch annotations
			},
		}
		accounts = append(accounts, account)
	}
	c.Negotiate(http.StatusOK, gin.Negotiate{Offered: ginutil.APIOfferedContentTypes, Data: accounts})
}

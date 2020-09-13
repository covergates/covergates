package user

import (
	"strconv"
	"time"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
)

// Token for API
type Token struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// HandleCreateToken for user
// @Summary create OAuth token
// @Tags User
// @Param name formData string false "token name"
// @Success 200 {object} string access token
// @Router /user/tokens [post]
func HandleCreateToken(service core.OAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.String(401, "")
			return
		}
		tokenName := c.PostForm("name")
		ctx := service.WithUser(c.Request.Context(), user)
		token, err := service.CreateToken(ctx, tokenName)
		if err != nil {
			c.Error(err)
			c.String(500, "")
			return
		}
		c.String(200, token.Access)
	}
}

// HandleListTokens for user
// @Summary list OAuth tokens
// @Tags User
// @Success 200 {object} []Token "list of tokens"
// @Router /user/tokens [get]
func HandleListTokens(service core.OAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(401, []*Token{})
			return
		}
		ctx := service.WithUser(c.Request.Context(), user)
		tokens, err := service.ListTokens(ctx)
		if err != nil {
			c.Error(err)
			c.JSON(500, []*Token{})
			return
		}
		result := make([]*Token, len(tokens))
		for i, token := range tokens {
			result[i] = &Token{
				ID:        token.ID,
				Name:      token.Name,
				CreatedAt: token.CreatedAt,
			}
		}
		c.JSON(200, result)
	}
}

// HandleDeleteToken with token id
// @Summary delete token with id
// @Tags User
// @Param id path integer true "token id"
// @Success 200 {object} Token "deleted token"
// @Router /user/tokens/{id} [delete]
func HandleDeleteToken(service core.OAuthService, store core.OAuthStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(401, &Token{})
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.Error(err)
			c.JSON(400, &Token{})
			return
		}
		token, err := store.Find(&core.OAuthToken{ID: uint(id)})
		if err != nil {
			c.Error(err)
			c.JSON(500, &Token{})
			return
		}
		ctx := service.WithUser(c.Request.Context(), user)
		if err := service.DeleteToken(ctx, token); err != nil {
			c.Error(err)
			c.JSON(500, &Token{})
			return
		}
		c.JSON(200, &Token{ID: token.ID, Name: token.Name})
	}
}

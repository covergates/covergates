package user

import (
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
)

type User struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

func HandleCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement user create
	}
}

// HandleGet login user
// @Summary Get login user
// @Tags User
// @Success 200 {object} User "user"
// @Failure 404 {string} string "error"
// @Router /user [get]
func HandleGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.String(404, "user not found")
			return
		}
		u := User{
			Login: user.Login,
			Email: user.Email,
		}
		c.JSON(200, u)
	}
}

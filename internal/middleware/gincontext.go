package middleware

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
)

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	c := ctx.Value("GinContextKey")
	_, err := GetUID(c.(*gin.Context))
	if err != nil {
		// appG.ResponseError(http.StatusUnauthorized, err, nil)
		// return

		return nil, fmt.Errorf("Access denied")
	}

	// or let it pass through
	return next(ctx)
}

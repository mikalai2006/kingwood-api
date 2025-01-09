package v1

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoles           = "roles"
	maxDistance         = "maxDistance"
	uid                 = "uid"
	authCtx             = "Auth"
)

func (h *HandlerV1) SetUserFromRequest(c *gin.Context) {
	appG := app.Gin{C: c}

	// fmt.Println(c.Request.Header)
	header := c.GetHeader(authorizationHeader)

	// socket auth
	if header == "" {
		header = c.GetHeader("token")
	}
	if header == "" {
		headerSecWebsocket := c.GetHeader("Sec-Websocket-Protocol")
		fmt.Println("headerSecWebsocket header1: ", header)
		headerSecWebsocketArray := strings.Split(headerSecWebsocket, ",")
		for i, _ := range headerSecWebsocketArray {
			if strings.Contains(headerSecWebsocketArray[i], "Bearer-") {
				header = strings.TrimSpace(strings.Replace(headerSecWebsocketArray[i], "Bearer-", "Bearer ", -1))
				fmt.Println("headerSecWebsocket header2: ", header)
			}
		}

	}

	// fmt.Println("header=", header)
	// jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	// fmt.Println("jwtCookie=", jwtCookie)

	if header == "" {
		// c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("empty auth header"))
		appG.ResponseError(http.StatusUnauthorized, errors.New("empty auth header"), nil)
		return
	}

	headerParts := strings.Split(header, " ")

	countParts := 2
	if len(headerParts) != countParts {
		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
		appG.ResponseError(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
		return
	}

	if headerParts[1] == "" {
		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
		appG.ResponseError(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
		return
	}

	// parse token
	// userId, err := h.services.Authorization.ParseToken(headerParts[1])
	// if err != nil {
	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	// 	return
	// }
	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	claims, err := tokenManager.Parse(string(headerParts[1]))
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	authData, err := h.Services.Authorization.GetAuth(claims.Subject)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	// fmt.Println("header2: ", len(headerParts))

	c.Set(userCtx, claims.Subject)
	c.Set(userRoles, claims.Roles)
	c.Set(maxDistance, claims.Md)
	c.Set(uid, claims.Uid)
	c.Set(authCtx, authData)
	// fmt.Println("claims.Uid=", claims.Uid)
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	// if user == nil {
	// 	// Abort the request with the appropriate error code
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// logrus.Printf("user session= %s", user)
	// // Continue down the chain to handler etc
	// c.Next()

	// c.JSON(http.StatusOK, Like)
	// fmt.Println("Set user from request", claims.Subject)
}

package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoles           = "roles"
	maxDistance         = "maxDistance"
	uid                 = "uid"
	authCtx             = "Auth"
)

// func SetUserIdentity(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	header := c.GetHeader(authorizationHeader)
// 	// fmt.Println("header=", header)
// 	// jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
// 	// fmt.Println("jwtCookie=", jwtCookie)

// 	if header == "" {
// 		// c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("empty auth header"))
// 		appG.ResponseError(http.StatusUnauthorized, errors.New("empty auth header"), nil)
// 		return
// 	}

// 	headerParts := strings.Split(header, " ")
// 	countParts := 2
// 	if len(headerParts) != countParts {
// 		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
// 		appG.ResponseError(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
// 		return
// 	}

// 	if headerParts[1] == "" {
// 		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
// 		appG.ResponseError(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
// 		return
// 	}

// 	// parse token
// 	// userId, err := h.services.Authorization.ParseToken(headerParts[1])
// 	// if err != nil {
// 	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 	// 	return
// 	// }
// 	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, nil)
// 		return
// 	}

// 	claims, err := tokenManager.Parse(string(headerParts[1]))
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, nil)
// 		return
// 	}
// 	c.Set(userCtx, claims.Subject)
// 	c.Set(userRoles, claims.Roles)
// 	c.Set(maxDistance, claims.Md)
// 	c.Set(uid, claims.Uid)
// 	// session := sessions.Default(c)
// 	// user := session.Get(userkey)
// 	// if user == nil {
// 	// 	// Abort the request with the appropriate error code
// 	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 	// 	return
// 	// }
// 	// logrus.Printf("user session= %s", user)
// 	// // Continue down the chain to handler etc
// 	// c.Next()
// }

// func SetUserIdentityGraphql(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	header := c.GetHeader(authorizationHeader)
// 	// fmt.Println("header=", header)
// 	// jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
// 	// fmt.Println("jwtCookie=", jwtCookie)

// 	authError := false

// 	if header == "" {
// 		// // c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("empty auth header"))
// 		// appG.ResponseError(http.StatusUnauthorized, errors.New("empty auth header"), nil)
// 		// return
// 		authError = true
// 	}

// 	headerParts := strings.Split(header, " ")
// 	countParts := 2
// 	if len(headerParts) != countParts || headerParts[1] == "" {
// 		// // c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
// 		// appG.ResponseError(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
// 		// return
// 		authError = true
// 	}

// 	if !authError {
// 		tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
// 		if err != nil {
// 			// c.AbortWithError(http.StatusUnauthorized, err)
// 			appG.ResponseError(http.StatusUnauthorized, err, nil)
// 			return
// 		}

// 		claims, err := tokenManager.Parse(string(headerParts[1]))
// 		if err != nil {
// 			// c.AbortWithError(http.StatusUnauthorized, err)
// 			appG.ResponseError(http.StatusUnauthorized, err, nil)
// 			return
// 		}
// 		c.Set(userCtx, claims.Subject)
// 		c.Set(userRoles, claims.Roles)
// 		c.Set(uid, claims.Uid)
// 		c.Set(maxDistance, claims.Md)
// 	} else {
// 		c.Set(userCtx, nil)
// 		c.Set(userRoles, nil)
// 		c.Set(uid, nil)
// 		c.Set(maxDistance, nil)

// 	}

// 	// id, ok := GetUserID(c)
// 	// if ok != nil {
// 	// 	fmt.Println("SetUserIdentityGraphql::: Not auth")
// 	// } else {
// 	// 	fmt.Println("SetUserIdentityGraphql::: AuthID=", id)
// 	// }
// }

func GetUserID(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user not found")
	}

	idInt, ok := id.(string)
	if !ok {
		return "", errors.New("user not found2")
	}

	return idInt, nil
}

func GetRoles(c *gin.Context) ([]string, error) {
	roles, ok := c.Get(userRoles)
	if !ok {
		return nil, errors.New("roles not found")
	}
	return roles.([]string), nil
}

// func GetMaxDistance(c *gin.Context) (int, error) {
// 	md, ok := c.Get(maxDistance)
// 	if !ok && md != nil {
// 		return 0, errors.New("max distance not found")
// 	}
// 	if md == nil {
// 		md = 5000000
// 	}
// 	return md.(int), nil
// }

func GetUID(c *gin.Context) (string, error) {
	id, ok := c.Get(uid)
	if !ok || id == "" {
		return "", errors.New("UID not found")
	}

	idString, ok := id.(string)
	if !ok {
		return "", errors.New("UID not found2")
	}

	return idString, nil
}

// func GetAuthFromCtx(c *gin.Context) (domain.Auth, error) {
// 	value, ex := c.Get(authCtx)
// 	if !ex {
// 		return domain.Auth{}, errors.New("auth is missing from ctx")
// 	}

// 	auth, ok := value.(domain.Auth)
// 	if !ok {
// 		return domain.Auth{}, errors.New("failed to convert value from ctx to domain.Auth")
// 	}

// 	return auth, nil
// }

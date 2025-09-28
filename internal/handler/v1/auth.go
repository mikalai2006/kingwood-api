package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *HandlerV1) registerAuth(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.POST("/sign-up", h.SignUp)
	auth.POST("/reset-password/:id", h.SetUserFromRequest, h.ResetPassword)
	auth.POST("/sign-in", h.SignIn)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.tokenRefresh)
	auth.GET("/refresh", h.tokenRefresh)
	auth.PATCH("/:id", h.SetUserFromRequest, h.UpdateAuth)
	auth.GET("/verification/:code", h.SetUserFromRequest, h.VerificationAuth)
	auth.GET("/iam", h.SetUserFromRequest, h.getIam)
}

// func (h *HandlerV1) updateAuth(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		appG.ResponseError(http.StatusUnauthorized, err, nil)
// 		return
// 	}

// 	var input domain.AuthInput
// 	if er := c.Bind(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	auth, err := h.Services.Authorization.UpdateAuth(userID, &input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, auth)
// }

func (h *HandlerV1) getIam(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	// TODO get token from body data.
	// var input *domain.RefreshInput

	// if err := c.BindJSON(&input); err != nil {
	// 	appG.Response(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// fmt.Println("ID=", userID)

	users, err := h.Services.User.Iam(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// get auth data for user
	// authData, err := h.Services.GetAuth(users.UserID.Hex())
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if !authData.ID.IsZero() {
	// 	// users.Md = authData.MaxDistance
	// 	users.Role = authData.RoleObject
	// 	fmt.Println("authData", authData)
	// }

	// // implementation max distance.
	// md, err := middleware.GetMaxDistance(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// users.Md = md

	// // implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// users.Roles = roles

	c.JSON(http.StatusOK, users)
}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.Auth true "account info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-up [post].
func (h *HandlerV1) SignUp(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("lang")
	if lang == "" {
		lang = h.i18n.Default
	}

	var input *domain.AuthInput
	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	input.Strategy = "local"

	// Check exist auth
	existAuth, err := h.Services.Authorization.ExistAuth(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if !existAuth.ID.IsZero() {
		appG.ResponseError(http.StatusBadRequest, errors.New("exist account"), nil)
		return
	}

	id, err := h.Services.Authorization.CreateAuth(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// create default
	// avatar := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon", id)

	newUser := domain.User{
		// Avatar: avatar,
		UserID: primitiveID,
		// Login:  input.Login,
		Name:     input.Name,
		Phone:    input.Phone,
		Hidden:   0,
		RoleId:   input.RoleId,
		PostId:   input.PostId,
		TypeWork: input.TypeWork,
		Oklad:    input.Oklad,
		TypePay:  input.TypePay,
		Birthday: &input.Birthday,
		// Post:  input.Post,
		// Roles: []string{"user"},
		// Md:    1,
	}
	document, err := h.Services.User.CreateUser(id, &newUser)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) ResetPassword(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("lang")
	if lang == "" {
		lang = h.i18n.Default
	}

	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.ResetPassword
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	// implementation roles for user.
	userForAuth, err := h.Services.User.GetUser(userID)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	fmt.Println("Reset password: ", input, id, userID, userForAuth.UserID.Hex())
	if userForAuth.UserID.Hex() != id {
		if !slices.Contains(roles, "auth-resetpass") {
			appG.ResponseError(http.StatusUnauthorized, domain.ErrNotRole, nil)
			return
		}
	}

	newPassword, err := h.Services.Authorization.ResetPassword(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, newPassword)
}

// @Summary SignIn
// @Tags auth
// @Description Login user
// @ID signin-account
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-in [post].
func (h *HandlerV1) SignIn(c *gin.Context) {
	appG := app.Gin{C: c}
	// jwt_cookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	// fmt.Println("+++++++++++++")
	// fmt.Printf("%s = %s",h.auth.NameCookieRefresh, jwt_cookie)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input *domain.AuthInput

	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if input.Strategy == "" {
		input.Strategy = "local"
	}

	if input.Email == "" && input.Login == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("error.notLogin"), nil)
		return
	}
	// if input.Password == "" {
	// 	appG.ResponseError(http.StatusBadRequest, errors.New("empty password"), nil)
	// 	return
	// }

	if input.Strategy == "local" {
		tokens, err := h.Services.Authorization.SignIn(input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		c.SetCookie(h.auth.NameCookieRefresh, tokens.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)

		c.JSON(http.StatusOK, domain.ResponseTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    tokens.ExpiresIn,
			ExpiresInR:   tokens.ExpiresInR,
		})
	}
	// else {
	// 	fmt.Print("JWT auth")
	// }
	// session.Set(userkey, input.Username)
	// session.Save()
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body domain.RefreshInput true "sign up info"
// @Success 200 {object} domain.ResponseTokens
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /users/auth/refresh [post].
func (h *HandlerV1) tokenRefresh(c *gin.Context) {
	appG := app.Gin{C: c}
	jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	// fmt.Sprintf("refresh Cookie %s = %s", h.auth.NameCookieRefresh, jwtCookie)
	// cookie_header := c.GetHeader("cookie")
	// fmt.Println("refresh Cookie_header = ", cookie_header)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input domain.RefreshInput

	// if jwtCookie == "" {
	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	// } else {
	// 	input.Token = jwtCookie
	// }
	// fmt.Println("refresh input.Token  = ", input.Token)
	// fmt.Println("jwtCookie  = ", jwtCookie)
	if input.Token == "" {
		input.Token = jwtCookie
	}

	if input.Token == "" && jwtCookie == "" {
		appG.ResponseError(http.StatusUnauthorized, errors.New("not found token"), nil)
		// c.JSON(http.StatusOK, gin.H{})
		// c.AbortWithStatus(http.StatusOK)
		return
	}

	fmt.Println("Refresh token=", input.Token)

	res, err := h.Services.Authorization.RefreshTokens(input.Token)
	if err != nil && err != mongo.ErrNoDocuments {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.SetCookie(h.auth.NameCookieRefresh, "", -1, "/", c.Request.URL.Hostname(), false, true)
		// res.RefreshToken = "expired"
		// res.AccessToken = "expired"
	} else {
		c.SetCookie(h.auth.NameCookieRefresh, res.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)
	}

	// userData, err := h.services.User.FindUser(domain.RequestParams{Filter: bson.D{{"userId": res.}}})
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// c.SetCookie(h.auth.NameCookieRefresh, res.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, domain.ResponseTokens{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
		ExpiresInR:   res.ExpiresInR,
	})
}

func (h *HandlerV1) Logout(c *gin.Context) {
	// session := sessions.Default(c)
	// session.Delete(userkey)
	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	// 	return
	// }
	appG := app.Gin{C: c}

	var input domain.RefreshInput

	jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	if jwtCookie == "" {
		if err := c.BindJSON(&input); err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
	} else {
		input.Token = jwtCookie
	}

	if input.Token == "" && jwtCookie == "" {
		c.JSON(http.StatusOK, gin.H{})
		c.AbortWithStatus(http.StatusOK)
		return
	}

	_, err := h.Services.Authorization.RemoveRefreshTokens(input.Token)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.SetCookie(h.auth.NameCookieRefresh, "", -1, "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *HandlerV1) VerificationAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	code := c.Param("code")
	if code == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("code empty"), nil)
		return
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if er := h.Services.Authorization.VerificationCode(userID, code); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *HandlerV1) UpdateAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }
	id := c.Param("id")

	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[domain.AuthInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	_, err := h.Services.Authorization.UpdateAuth(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

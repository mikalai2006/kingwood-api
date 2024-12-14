package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/service"
)

func (h *HandlerV1) registerWs(router *gin.RouterGroup) {
	wsrouter := router.Group("/ws")
	wsrouter.GET("/:roomId", h.goWs)
}

func (h *HandlerV1) goWs(c *gin.Context) {
	// appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }
	// d, _ := c.Get(authCtx)
	// fmt.Println(userID, d)
	roomId := c.Param("roomId")
	service.ServeWS(c, roomId, &h.hub, h.Services)

	// c.JSON(http.StatusOK, nil)
	c.Done()
}

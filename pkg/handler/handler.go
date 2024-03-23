package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"testApi/pkg/service"

	_ "github.com/MrPrograMax/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/group", h.createGroup)
		api.GET("/groups", h.getAllGroups)
		api.GET("/group/:id", h.getGroupById)
		api.PUT("/group/:id", h.updateGroup)
		api.DELETE("/group/:id", h.deleteGroup)

		items := api.Group("/group")
		{
			items.POST("/:id/participant", h.createParticipant)
			items.POST("/:id/toss", h.tossParticipant)
			items.DELETE("/:id/participant/:participantId", h.deleteParticipant)
			items.GET("/:id/participant/:participantId/recipient", h.getInfoAboutRecipient)
		}
	}

	return router
}

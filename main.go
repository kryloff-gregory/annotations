package main

import (
	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"main/annotation"
	"main/auth"
	"main/controller"
	"main/manager"
	"main/services"
	"main/user"
	"main/video"
)

func main() {
	conn, err := leveldb.OpenFile("db/annotations", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	lockService := services.NewLockService()

	videoDao := video.NewDAO(conn)
	videoManager := manager.NewVideoManager(videoDao, lockService)
	videoController := controller.NewVideoController(videoManager)

	annDAO := annotation.NewDAO(conn)
	annManager := manager.NewAnnotationManager(videoDao, annDAO, lockService)
	annController := controller.NewAnnotationController(annManager)

	userProvider, err := user.NewProvider("users.json")
	if err != nil {
		log.Fatal(err)
	}
	loginController := controller.NewLoginController(userProvider)

	middleware := auth.NewMiddleware(userProvider)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/login", loginController.Login)

	v1.Use(middleware.Authenticate())
	{
		v1.GET("/video", videoController.GetVideo)
		v1.PUT("/video", videoController.CreateVideo)
		v1.DELETE("/video", videoController.DeleteVideo)
		v1.GET("/annotation", annController.GetAnnotationsForVideo)
		v1.PUT("/annotation", annController.CreateAnnotation)
		v1.POST("/annotation", annController.UpdateAnnotation)
		v1.DELETE("/annotation", annController.DeleteAnnotation)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.Run(":5001")

}

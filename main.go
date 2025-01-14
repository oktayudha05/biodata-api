package main

import (
	"biodata-server/controller"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	mahasiswa := router.Group("/mahasiswa")
	{
		mahasiswa.GET("/", controller.GetMhs)
	}

	router.Run(":3000")
}
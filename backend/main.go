package main

import (
	"github.com/Marcus-Nastasi/docx2pdf/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.POST("/convert", controller.Convert)
	server.Run(":8081")
}

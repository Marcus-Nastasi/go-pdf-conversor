package main

import (
	"github.com/Marcus-Nastasi/docx2pdf/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.POST("/convert-on-machine", controller.ConvertOnMachine)
	server.POST("/convert-upload", controller.ConvertUpload)
	server.Run(":8081")
}

package main

import (
	"net/http"

	"github.com/Marcus-Nastasi/docx2pdf/application"
	"github.com/Marcus-Nastasi/docx2pdf/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	var converter application.Converter = application.NewLibreOfficeConverter(20)
	controller := controller.NewConverterController(&converter)

	server.LoadHTMLFiles("../frontend/index.html")
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})
	server.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"status": "all clear"})
	})
	server.POST("/convert", controller.ConvertUpload)
	// server.POST("/convert-on-machine", controller.ConvertOnMachine)

	server.Run(":8081")
}

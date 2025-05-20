package controller

import (
	"net/http"

	"github.com/Marcus-Nastasi/docx2pdf/application"
	"github.com/gin-gonic/gin"
)

type Path struct {
	Path string `json:"path"`
}

func Convert(ctx *gin.Context) {
	var docxPath Path
	err := ctx.BindJSON(&docxPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
	}
	var converter application.Convert = &application.LibreOfficeConverter{}
	pdfPath, err := converter.ConvertToPdf(docxPath.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	ctx.JSON(http.StatusOK, map[string]string{"path": pdfPath})
}

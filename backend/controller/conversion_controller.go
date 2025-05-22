package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Marcus-Nastasi/docx2pdf/application"
	"github.com/gin-gonic/gin"
)

type Path struct {
	Path string `json:"path"`
}

func ConvertOnMachine(ctx *gin.Context) {
	var docxPath Path
	err := ctx.BindJSON(&docxPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	var converter application.Convert = &application.LibreOfficeConverter{}
	pdfPath, err := converter.LocalConvertToPdf(docxPath.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"path": pdfPath})
}

func ConvertUpload(ctx *gin.Context) {
	// Recover file
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Error reading file: %s", err.Error())
		return
	}
	// Garantee file is closed after
	defer file.Close()
	// Clear blank spaces on name
	fileName := strings.ReplaceAll(header.Filename, " ", "_")
	// Init converter
	var converter application.Convert = &application.LibreOfficeConverter{}
	// Saves, convert the file, than clean tmp dir
	pdfBytes, err := converter.ConvertFromUpload(&fileName, &file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	// Define response headers for download
	ctx.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", converter.ChangeExtension(header.Filename, "pdf")),
	)
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

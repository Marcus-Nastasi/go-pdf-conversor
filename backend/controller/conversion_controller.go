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

type ConverterController struct {
	Converter application.Converter
}

func NewConverterController(c *application.Converter) *ConverterController {
	return &ConverterController{
		Converter: *c,
	}
}

func (c *ConverterController) ConvertOnMachine(ctx *gin.Context) {
	var docxPath Path
	err := ctx.BindJSON(&docxPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	pdfPath, err := c.Converter.LocalConvertToPdf(&docxPath.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"path": pdfPath})
}

func (c *ConverterController) ConvertUpload(ctx *gin.Context) {
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
	// Saves, convert the file, than clean tmp dir
	pdfBytes, err := c.Converter.ConvertFromUpload(&fileName, &file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	ext := "pdf"
	// Define response headers for download
	ctx.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", c.Converter.ChangeExtension(&header.Filename, &ext)),
	)
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

package controller

import (
	"archive/zip"
	"bytes"
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
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"path": pdfPath})
}

func (c *ConverterController) ConvertUpload(ctx *gin.Context) {
	ext := "pdf"
	// Recover file
	mtp, err := ctx.MultipartForm()
	files := mtp.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "no files to convert"})
		return
	}

	// Run the ConvertFromUpload for 1 file input
	if len(files) < 2 {
		file := files[0]
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		// Clear blank spaces on name
		fileName := strings.ReplaceAll(file.Filename, " ", "_")
		openedFile, err := file.Open()
		// Garantee file is closed right after
		defer openedFile.Close()
		// Saves, convert the file, than clean tmp dir
		pdfBytes, err := c.Converter.ConvertFromUpload(&fileName, openedFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		// Define response headers for download
		ctx.Header(
			"Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", c.Converter.ChangeExtension(&file.Filename, &ext)),
		)
		ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
		return
	}

	// Runs the ConvertMultiple for multiple file inputs
	if len(files) > 1 {
		var converted_files [][]byte
		converted_files, err := c.Converter.ConvertMultiple(files)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		for i, pdfBytes := range converted_files {
			originalFileName := files[i].Filename
			pdfFileName := c.Converter.ChangeExtension(&originalFileName, &ext)
			f, err := zipWriter.Create(pdfFileName)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			_, err = f.Write(pdfBytes)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		
		err = zipWriter.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Send zip file
		ctx.Header("Content-Disposition", "attachment; filename=\"converted_pdfs.zip\"")
		ctx.Data(http.StatusOK, "application/zip", buf.Bytes())
		return
	}
}

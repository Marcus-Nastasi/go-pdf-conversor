package service

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Convert interface {
	ConvertToPdf(docxPath string) (string, error)
}

type LibreOfficeConverter struct{}

func (l LibreOfficeConverter) ConvertToPdf(docxPath string) (string, error) {
	arg0 := "lowriter"
	arg1 := "--invisible" //This command is optional, it will help to disable the splash screen of LibreOffice.
	arg2 := "--convert-to"
	arg3 := "pdf:writer_pdf_Export"
	outputDir := filepath.Dir(docxPath)
	_, err := exec.Command(arg0, arg1, arg2, arg3, "--outdir", outputDir, docxPath).Output()
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
		return "", err
	}
	return outputDir, nil
}

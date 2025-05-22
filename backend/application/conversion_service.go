package application

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

type Convert interface {
	LocalConvertToPdf(docxPath string) (string, error)
	ConvertFromUpload(fileName *string, file *multipart.File) ([]byte, error)
	ChangeExtension(filename, newExt string) string
}

type LibreOfficeConverter struct{}

func (l LibreOfficeConverter) LocalConvertToPdf(docxPath string) (string, error) {
	arg0 := "lowriter"
	arg1 := "--invisible" // disable the splash screen of LibreOffice.
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

func (l LibreOfficeConverter) ConvertFromUpload(fileName *string, file *multipart.File) ([]byte, error) {
	// Creates temporary file
	tempDir := os.TempDir()
	tempInputFile := filepath.Join(tempDir, "upload_"+*fileName)
	// Do the upload of original file on tmp dir
	outFile, err := os.Create(tempInputFile)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()
	// Copies all the original file content to the temporary file created
	io.Copy(outFile, *file)
	// Convert file with LibreOffice
	convertedDir, err := l.LocalConvertToPdf(tempInputFile)
	if err != nil {
		return nil, err
	}
	// Define the generated pdf path
	outputPDF := filepath.Join(convertedDir, l.ChangeExtension("upload_"+*fileName, "pdf"))
	// Reads the pdf
	pdfBytes, err := os.ReadFile(outputPDF)
	if err != nil {
		return nil, err
	}
	// Clean temporary files
	os.Remove(tempInputFile)
	os.Remove(outputPDF)
	return pdfBytes, nil
}

// Changes the file extension
func (l LibreOfficeConverter) ChangeExtension(filename, newExt string) string {
	return fmt.Sprintf("%s.%s", filename[:len(filename)-len(filepath.Ext(filename))], newExt)
}

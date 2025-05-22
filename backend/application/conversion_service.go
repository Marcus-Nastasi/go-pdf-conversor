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
	ConvertFromUpload(fileName string, file multipart.File) ([]byte, error)
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

func (l LibreOfficeConverter) ConvertFromUpload(fileName string, file multipart.File) ([]byte, error) {
	// Creates temporary file
	tempDir := os.TempDir()
	ext := filepath.Ext(fileName)
	if ext == "" {
		return nil, &exec.Error{}
	}
	tempInputFile := filepath.Join(tempDir, "upload_"+fileName)
	outFile, err := os.Create(tempInputFile)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()
	io.Copy(outFile, file)
	// 3. Converte o arquivo com LibreOffice
	convertedDir, err := l.LocalConvertToPdf(tempInputFile)
	if err != nil {
		return nil, err
	}
	// 4. Define o caminho do PDF gerado
	outputPDF := filepath.Join(convertedDir, changeExtension("upload_"+fileName, "pdf"))
	// 5. Lê o PDF gerado
	pdfBytes, err := os.ReadFile(outputPDF)
	if err != nil {
		return nil, err
	}
	// 7. Limpa arquivos temporários
	os.Remove(tempInputFile)
	os.Remove(outputPDF)
	return pdfBytes, nil
}

// changeExtension muda a extensão de um arquivo
func changeExtension(filename, newExt string) string {
	return fmt.Sprintf("%s.%s", filename[:len(filename)-len(filepath.Ext(filename))], newExt)
}

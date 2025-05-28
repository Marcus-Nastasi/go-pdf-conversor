package application

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Converter interface {
	LocalConvertToPdf(docxPath *string) (string, error)
	ConvertFromUpload(fileName *string, file multipart.File) ([]byte, error)
	ConvertMultiple(fileHeaders []*multipart.FileHeader) ([][]byte, error)
	ChangeExtension(filename, newExt *string) string
}

type LibreOfficeConverter struct{
	libreOfficePool chan struct{}
}

func NewLibreOfficeConverter(n int) *LibreOfficeConverter {
	return &LibreOfficeConverter{
		libreOfficePool: make(chan struct{}, n),
	}
}

func (l *LibreOfficeConverter) LocalConvertToPdf(docxPath *string) (string, error) {
	// Acquire a slot in the pool
	l.libreOfficePool <- struct{}{}
	defer func() { <- l.libreOfficePool }()
	// Define output and unique user profile directory
	outputDir := filepath.Dir(*docxPath)
	profileDir := filepath.Join(os.TempDir(), "lo-profile-"+uuid.New().String())
	err := os.MkdirAll(profileDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create user profile dir: %w", err)
	}
	// Build command
	cmd := exec.Command(
		"lowriter",
		"--invisible",
		"--convert-to", "pdf:writer_pdf_Export",
		fmt.Sprintf("-env:UserInstallation=file://%s", profileDir),
		"--outdir", outputDir,
		*docxPath,
	)
	// Capture stderr for diagnostics
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stderr // Optional: capture all output
	err = cmd.Run()
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", stderr.String())
		return "", fmt.Errorf("LibreOffice failed: %s", stderr.String())
	}
	return outputDir, nil
}

func (l *LibreOfficeConverter) ConvertFromUpload(fileName *string, file multipart.File) ([]byte, error) {
	// Creates temporary file
	tempDir := os.TempDir()
	newUuid := uuid.New()
	tempInputFile := filepath.Join(tempDir, "upload_" + newUuid.String() + "_" + *fileName)
	// Do the upload of original file on tmp dir
	outFile, err := os.Create(tempInputFile)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()
	// Copies all the original file content to the temporary file created
	io.Copy(outFile, file)
	// Convert file with LibreOffice
	convertedDir, err := l.LocalConvertToPdf(&tempInputFile)
	if err != nil {
		return nil, err
	}
	ext := "pdf"
	uploadFileName := "upload_" + newUuid.String() + "_" + *fileName
	// Define the generated pdf path
	outputPDF := filepath.Join(convertedDir, l.ChangeExtension(&uploadFileName, &ext))
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

func (l *LibreOfficeConverter) ConvertMultiple(fileHeaders []*multipart.FileHeader) ([][]byte, error) {
	if fileHeaders == nil {
		return nil, errors.New("no files to convert")
	}
	var (
		pdfBytes = make([][]byte, len(fileHeaders))
		wg       sync.WaitGroup 	// Used to wait for all goroutines to finish
		mu       sync.Mutex     	// Used to avoid data races when writing to shared slice
		errOnce  sync.Once      	// Ensures only one error is captured
		retErr   error          	// Holds the first error found
	)
	for i, fileHeader := range fileHeaders {
		wg.Add(1)
		go func(i int, fh *multipart.FileHeader) {
			defer wg.Done()
			file, err := fileHeader.Open()
			if err != nil {
				errOnce.Do(func() { retErr = err })
				return
			}
			defer file.Close() 
			fileName := strings.ReplaceAll(fh.Filename, " ", "_")
			pdfByte, err := l.ConvertFromUpload(&fileName, file)
			if err != nil {
				errOnce.Do(func() { retErr = err })
				return
			}
			mu.Lock()
			pdfBytes[i] = pdfByte
			mu.Unlock()
		}(i, fileHeader)
	}
	wg.Wait()
	if retErr != nil {
		return nil, retErr
	}
	return pdfBytes, nil
}

// Changes the file extension
func (l *LibreOfficeConverter) ChangeExtension(filename, newExt *string) string {
	return strings.Replace(*filename, filepath.Ext(*filename), "." + *newExt, 1)
}

// func (l *LibreOfficeConverter) LocalConvertToPdf(docxPath *string) (string, error) {
// 	// l.libreOfficeMutex.Lock()
// 	// Acquire a slot in the pool
// 	l.libreOfficePool <- struct{}{}
// 	defer func() { <- l.libreOfficePool }()

// 	arg0 := "lowriter"
// 	arg1 := "--invisible" // disable the splash screen of LibreOffice.
// 	arg2 := "--convert-to"
// 	arg3 := "pdf:writer_pdf_Export"
// 	outputDir := filepath.Dir(*docxPath)
// 	_, err := exec.Command(arg0, arg1, arg2, arg3, "--outdir", outputDir, *docxPath).Output()
// 	if err != nil {
// 		fmt.Printf("[ERROR]: %s\n", err.Error())
// 		return "", err
// 	}
// 	// l.libreOfficeMutex.Unlock()
// 	return outputDir, nil
// }

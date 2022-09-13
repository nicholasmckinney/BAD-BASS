package resource

import (
	"archive/zip"
	"bytes"
	"crypto/rc4"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Archive struct {
	Filepath string
}

func (ar *Archive) Unpack(outputDir string) error {
	if ar.Filepath == "" {
		return fmt.Errorf("no input filepath given")
	}

	inputPath, err := filepath.Abs(ar.Filepath)
	if err != nil {
		return fmt.Errorf("unable to get path for file (%s): %w", ar.Filepath, err)
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("unable to read file (%s): %w", inputPath, err)
	}

	zipBytes := make([]byte, len(inputData))
	cipher, err := rc4.NewCipher([]byte("zazalul"))
	cipher.XORKeyStream(zipBytes, inputData)

	zipReader := bytes.NewReader(zipBytes)
	reader, err := zip.NewReader(zipReader, int64(zipReader.Len()))

	for _, file := range reader.File {
		outputPath := filepath.Join(outputDir, file.Name)

		if file.FileInfo().IsDir() {
			fmt.Printf("[+] Processing directory: %s\n", file.Name)
			err = os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error while processing directory in zip (%s): %w", file.Name, err)
			}
			continue
		}

		fmt.Printf("[+] Processing file: %s\n", file.Name)
		err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error while processing file in zip (%s): %w", file.Name, err)
		}

		dstFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("error while opening output file (%s): %w", outputPath, err)
		}

		archivedFileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error while opening archived file (%s): %w", file.Name, err)
		}

		_, err = io.Copy(dstFile, archivedFileReader)
		if err != nil {
			return fmt.Errorf("error while writing content to output file (%s) from archived file (%s): %w", outputPath, file.Name, err)
		}
		dstFile.Close()
		archivedFileReader.Close()
	}
	return nil
}

package resource

import (
	"PhishRod/internal"
	"archive/zip"
	"bytes"
	"crypto/rc4"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Builder struct {
	directories []directory
	matcher     []byte
}

type directory struct {
	Path  string
	Key   string
	Files map[string]bytes.Buffer
}

//type configuration struct {
//	Entries []configurationEntry `xml:"entry"`
//}

type configurationEntry struct {
	Key  string
	Path string
}

func (r *Builder) AddDirectory(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("unable to open path: %s", path)
	}
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file metadata: %s", path)
	}

	if !info.IsDir() {
		return fmt.Errorf("specified path is not a directory: %s", path)
	}
	d := directory{
		Path: filepath.Base(path),
		Key:  internal.RandomString(8),
	}
	r.directories = append(r.directories, d)
	return nil
}

func (r *Builder) AddMatchFile(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("unable to open path: %s", path)
	}

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file metadata: %s", path)
	}

	if info.IsDir() {
		return fmt.Errorf("specified path is a directory: %s", path)
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read match file (%s): %w", path, err)
	}

	r.matcher = fileContent
	return nil
}

func (r *Builder) Build() ([]byte, error) {
	var buf bytes.Buffer
	//var configuration configuration
	var entries []configurationEntry

	if r.matcher == nil {
		return nil, fmt.Errorf("match file has not been added")
	}

	outputWriter := zip.NewWriter(&buf)

	for _, directory := range r.directories {
		//cipher, err := rc4.NewCipher([]byte(directory.Key))
		//if err != nil {
		//	return nil, fmt.Errorf("unable to generate cipher for directory (%s): %w", directory.Path, err)
		//}

		//fmt.Printf(
		//	"[+] Encrypting files in (%s) with key (%s)\n",
		//	directory.Path, directory.Key,
		//)

		err := filepath.WalkDir(directory.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("[-] Skipping file (%s) due to error: %v\n", path, err)
				return nil
			}
			if d.IsDir() {
				return nil
			}
			fmt.Printf("[+] Packing file: %s\n", path)
			var fileContent []byte
			fileContent, err = os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("unable to read file (%s): %w", path, err)
			}
			//ciphertext := make([]byte, len(fileContent))
			//cipher.XORKeyStream(ciphertext, fileContent)

			writer, err := outputWriter.Create(path)
			if err != nil {
				return fmt.Errorf("unable to create archived file (%s): %w", path, err)
			}
			buf := bytes.NewBuffer(fileContent)
			if _, err := io.Copy(writer, buf); err != nil {
				return fmt.Errorf("unable to write to archived file (%s): %w", path, err)
			}

			entries = append(entries, configurationEntry{Path: path, Key: directory.Key})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error while walking directory (%s): %v", directory.Path, err)
		}

	}

	//configuration.Entries = entries
	//marshalled, err := xml.Marshal(configuration)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to marshal configuration xml: %w", err)
	//}
	//keyfile, err := outputWriter.Create("map")
	//if err != nil {
	//	return nil, fmt.Errorf("error while creating key file: %v", err)
	//}

	//_, err = io.Copy(keyfile, bytes.NewBuffer(marshalled))
	//if err != nil {
	//	return nil, fmt.Errorf("error while writing configuration to archive: %w", err)
	//}

	matchfile, err := outputWriter.Create("conf.xml")
	_, err = io.Copy(matchfile, bytes.NewBuffer(r.matcher))
	if err != nil {
		return nil, fmt.Errorf("error while writing match file to archive: %w", err)
	}

	outputWriter.Close()
	zipBytes := buf.Bytes()

	cipher, err := rc4.NewCipher([]byte("zazalul"))
	if err != nil {
		return nil, fmt.Errorf("error while creating cipher for zip: %w", err)
	}

	outputBytes := make([]byte, len(zipBytes))
	cipher.XORKeyStream(outputBytes, zipBytes)
	fmt.Println("[+] Done!")
	return outputBytes, nil
}

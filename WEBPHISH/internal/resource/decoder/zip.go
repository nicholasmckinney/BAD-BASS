package decoder

import (
	"Webphish/internal"
	"Webphish/internal/resource"
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rc4"
	"io"
)

type RC4ZipDecoder struct {
	Key string
}

func NewRC4ZipDecoder() resource.Decoder {
	return &RC4ZipDecoder{Key: resource.DefaultKey}
}

func (d *RC4ZipDecoder) Decode(data []byte) ([]resource.File, internal.ErrorCode) {
	zipBytes := make([]byte, len(data))
	cipher, err := rc4.NewCipher([]byte(resource.DefaultKey))
	var files []resource.File

	if err != nil {
		return nil, internal.ERR_DECODER_CIPHER_CREATION
	}
	cipher.XORKeyStream(zipBytes, data)

	bytesReader := bytes.NewReader(zipBytes)
	reader, err := zip.NewReader(bytesReader, int64(bytesReader.Len()))

	if err != nil {
		return nil, internal.ERR_DECODER_UNPACK
	}

	for _, file := range reader.File {
		var buf bytes.Buffer
		bytesWriter := bufio.NewWriter(&buf)
		f := resource.File{
			Path:        file.Name,
			IsDirectory: file.FileInfo().IsDir(),
		}

		reader, err := file.Open()
		if err != nil {
			return nil, internal.ERR_DECODER_READ_CONTENT
		}
		_, err = io.Copy(bytesWriter, reader)
		if err != nil {
			return nil, internal.ERR_DECODER_READ_CONTENT
		}
		f.Content = buf.Bytes()

		files = append(files, f)
	}

	return files, internal.GenericSuccess
}

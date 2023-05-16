package helper

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(filename string, file multipart.File) {
	dst, err := os.Create("public/" + filename)
	defer func() {
		err := dst.Close()
		if err != nil {
			PanicIfError(err)
		}
	}()
	PanicIfError(err)

	_, err = io.Copy(dst, file)
	PanicIfError(err)
}

func RemoveFile(filename string) {
	err := os.Remove("public/" + filename)
	PanicIfError(err)
}

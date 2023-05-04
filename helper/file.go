package helper

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(filename string, file multipart.File) {
	dst, err := os.Create("public/" + filename)
	PanicIfError(err)
	defer dst.Close()

	_, err = io.Copy(dst, file)
	PanicIfError(err)
}

func RemoveFile(filename string) {
	err := os.Remove("public/" + filename)
	PanicIfError(err)
}

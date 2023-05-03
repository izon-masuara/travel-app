package helper

import "os"

func RemoveFile(filename string) {
	err := os.Remove("public/" + filename)
	PanicIfError(err)
}

package addon

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	Path = "/vendor/assets/picture/"
)

func GetDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func Upload(collection, name string, file multipart.File) (path string, err error) {
	path = path + collection + "/" + name + ".png"
	out, err := os.Create(path)
	if err != nil {
		err = errors.New("error creating file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		err = errors.New("error while copying file")
		return
	}
	return
}

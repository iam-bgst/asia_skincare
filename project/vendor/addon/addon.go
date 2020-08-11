package addon

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	Path = GetDir() + "/vendor/assets/picture/"
)

func GetDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func Upload(collection, name string, file multipart.File) (path string, err error) {
	path = Path + collection
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	out, err := os.Create(path + "/" + name + ".png")
	if err != nil {
		log.Println(err)
	} // Err Handling
	defer out.Close()

	_, err = io.Copy(out, file)
	path = "/picture/" + collection + "/" + name + ".png"
	return
}

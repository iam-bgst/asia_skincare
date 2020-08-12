package addon

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var (
	Path = GetDir() + "/vendor/assets/picture/"
)

func GetDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func Upload(collection, name string, file *multipart.FileHeader, c *gin.Context) (path string, err error) {
	fmt.Println(c.Request)
	path = Path + collection
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	out, err := os.Create(path + "/" + name + ".png")
	if err != nil {
		log.Println(err)
	} // Err Handling
	defer out.Close()
	err = c.SaveUploadedFile(file, path)
	path = "/picture/" + collection + "/" + name + ".png"
	return
}

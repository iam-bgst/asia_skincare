package addon

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	Path = GetDir() + "/vendor/assets/picture/"
)

func GetDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func Upload(collection, name string, img string) (path string, err error) {
	if img == "" {
		path = "/picture/" + collection + "/" + name + ".png"
		return
	}
	path = Path + collection
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	// out, err := os.Create(path + "/" + name + ".png")
	idx := strings.Index(img, ";base64,")
	if idx < 0 {
		err = errors.New("error indexing image")
		return
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img[idx+8:]))
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(reader)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/"+name+".png", buff.Bytes(), 0644)
	if err != nil {
		return
	}
	path = "/picture/" + collection + "/" + name + ".png"

	return
}

var letterRune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomCode(length int, origin bool) string {
	if origin {
		b := make([]rune, length)
		for i := range b {
			b[i] = letterRune[rand.Intn(len(letterRune))]
		}
		return string(b)
	} else {
		letterRune = append(letterRune, []rune("0123456789")...)
		b := make([]rune, length)
		for i := range b {
			b[i] = letterRune[rand.Intn(len(letterRune))]
		}
		return string(b)
	}
}

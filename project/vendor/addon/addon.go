package addon

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maddevsio/fcm"
	"github.com/swaggo/swag"
)

var (
	Path = GetDir() + "/vendor/assets/picture/"
)

const (
	NORMAL = fcm.PriorityNormal
	HIGH   = fcm.PriorityHigh

	TRANSACTION = "transaction"
	POINT       = "point"
	REDEEM      = "redeem"
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
			rand.Seed(time.Now().UnixNano())
			b[i] = letterRune[rand.Intn(len(letterRune))]
		}
		return string(b)
	} else {
		letterRune = append(letterRune, []rune("0123456789")...)
		b := make([]rune, length)
		for i := range b {
			rand.Seed(time.Now().UnixNano())
			b[i] = letterRune[rand.Intn(len(letterRune))]
		}
		return string(b)
	}
}

func DateSameOrNot(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

type Data struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func PushNotif(token, priority string, data Data, content ...string) {
	c := fcm.NewFCM("AAAAxAsG0w8:APA91bFzjef54BSksgSaOUTf3hKk5dRmMuz-i_HgvCYEH6s8_FyyAqMldf71c-JFUNazCTiNXLE_E-lBsI9qM5uP2gEsnOpKdWy2QnIVQjihZPQyTQhCWGdW_s0sV0kjcRVbgHOqgeiL")
	// c := fcm.NewFCM("AAAAldYRR1Y:APA91bHKQZTWlqU-X2KvCEvzlT-ukpj-siHtEIkzaKXyhenrqVzODOItXsN27j8jE_Pz8J8I7stjrYIo6mY-GoyzipnvEEscyJeV1bmKdvWihkajvBPxWK4KTmSE7Cz_gFDjsjmK95tL")
	//token := "dx6yRgG-c_w:APA91bFfxc84LhJ1JWQORBEYujBYUDXd1IBSap4Zf8Z5jGq-xDH-enRTsMIazfVsMvCYp_uhfCIKjiMfr65BwP2X_i7mv-wLk5RRHHGyx_ilUeHnOsLiRouKxspZYqlL2bnKrG9N3lWj"
	var datafix gin.H
	if data.Type == TRANSACTION {
		datafix = gin.H{
			"type": data.Type,
			"content": gin.H{
				"navigate": content[0],
			},
		}
	} else if data.Type == REDEEM {
		datafix = gin.H{
			"type": data.Type,
			"content": gin.H{
				"navigate": content[0],
			},
		}
	} else if data.Type == POINT {
		datafix = gin.H{
			"type": data.Type,
			"content": gin.H{
				"navigate": content[0],
			},
		}
	}

	response, err := c.Send(fcm.Message{
		Data:             datafix,
		RegistrationIDs:  []string{token},
		ContentAvailable: true,
		Priority:         priority,
		Notification: fcm.Notification{
			Title: data.Title,
			Body:  data.Body,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Status Code   :", response.StatusCode)
	fmt.Println("Success       :", response.Success)
	fmt.Println("Fail          :", response.Fail)
	fmt.Println("Canonical_ids :", response.CanonicalIDs)
	fmt.Println("Topic MsgId   :", response.MsgID)
}

var (
	swagger = swag.Register("Asia Skincare")
)

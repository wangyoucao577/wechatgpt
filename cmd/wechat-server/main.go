package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var flags struct {
	apiKey string
	token  string
	port   uint
}

func init() {
	flag.StringVar(&flags.apiKey, "api_key", "", "Your api_key of OpenAI platform.")
	flag.StringVar(&flags.token, "token", "", "Your token to verify wechat platform requests.")
	flag.UintVar(&flags.port, "p", 8080, "Listen port.")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	token := flags.token

	r := gin.Default()
	r.GET("/wx", func(c *gin.Context) {
		// https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Getting_Started_Guide.html

		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")

		l := []string{token, timestamp, nonce}
		sort.Strings(l)
		hashcode := sha1.Sum([]byte(strings.Join(l, "")))
		hashcodeStr := hex.EncodeToString(hashcode[:])

		if hashcodeStr == signature {
			c.String(http.StatusOK, echostr)
			return
		}
		c.String(http.StatusBadRequest, "")

	})
	r.Run(":" + strconv.Itoa(int(flags.port))) // listen and serve on 0.0.0.0:port
}

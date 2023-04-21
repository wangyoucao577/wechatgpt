package main

import (
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var flags struct {
	port uint
}

func init() {
	flag.UintVar(&flags.port, "p", 8080, "Listen port.")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	r := gin.Default()
	r.GET("/wx", wxValidationHandler)
	r.POST("/wx", wxMessageHandler)
	r.Run(":" + strconv.Itoa(int(flags.port))) // listen and serve on 0.0.0.0:port
}

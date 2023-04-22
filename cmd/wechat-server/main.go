package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var flags struct {
	port uint
}

func init() {
	flag.UintVar(&flags.port, "p", 8080, "Listen port.")
}

var questionsChan chan question
var answersChan chan string

func main() {
	flag.Parse()
	defer glog.Flush()

	questionsChan = make(chan question, 10) // max cache 10 questions
	answersChan = make(chan string, 10)     // max cache 10 answers

	go func() { // chatgpt tasks
		for {
			q, ok := <-questionsChan
			if !ok {
				glog.Infoln("questions channel closed")
				return
			}

			answer := chatgpt(q.content, time.Minute)

			// TODO: answers chan per user
			// TODO: split answers for wechat
			answersChan <- answer
		}
	}()
	defer close(questionsChan)
	defer close(answersChan)

	r := gin.Default()
	r.GET("/wx", wxValidationHandler)
	r.POST("/wx", wxMessageHandler)
	r.Run(":" + strconv.Itoa(int(flags.port))) // listen and serve on 0.0.0.0:port
}

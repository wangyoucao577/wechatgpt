package main

import (
	"flag"
	"strconv"
	"sync"
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
var answersMap sync.Map // user -> answersChan

func main() {
	flag.Parse()
	defer glog.Flush()

	questionsChan = make(chan question, 10) // max cache 10 questions

	go func() { // chatgpt tasks
		for {
			q, ok := <-questionsChan
			if !ok {
				glog.Infoln("questions channel closed")
				return
			}

			answer := chatgpt(q.content, time.Minute)

			if _, ok := answersMap.Load(q.user); !ok {
				answersMap.Store(q.user, make(chan string, 10)) // max cache 10 answers)
			}
			answersChanAny, _ := answersMap.Load(q.user)
			answersChan := answersChanAny.(chan string)

			// TODO: split answers by length for wechat
			answersChan <- answer
		}
	}()
	defer close(questionsChan)
	defer answersMap.Range(func(k, v any) bool {
		close(v.(chan string))
		return true
	})

	r := gin.Default()
	r.GET("/wx", wxValidationHandler)
	r.POST("/wx", wxMessageHandler)
	glog.Errorln(r.Run(":" + strconv.Itoa(int(flags.port)))) // listen and serve on 0.0.0.0:port
}

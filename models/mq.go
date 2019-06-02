package models

import (
	"blogmesssage"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

/**
负责消息队列里的发送逻辑
*/
var (
	mqConn *amqp.Connection
	maCh   *amqp.Channel
	q      amqp.Queue
)

func init() {
	var err error
	mqConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		beego.BeeLogger.Error("%s", err)
	}

	maCh, err = mqConn.Channel()

	q, err = maCh.QueueDeclare(
		"article",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		beego.BeeLogger.Error("%s", err)
	}
}

func SendArticleMessage(message *blogmesssage.ArticleMessage) {
	content, err := json.Marshal(message)
	if err != nil {
		return
	}
	err = maCh.Publish("",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        content,
		},
	)

	if err != nil {
		log.Printf("%s", err.Error())
	}
}

//todo 添加移除操作

func getRecommendArticleIds(ids []int64) ([]int64, error) {
	content, err := json.Marshal(ids)
	q, err := maCh.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	corrId := randomString(32)
	err = maCh.Publish("rpc",
		"recommended",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          content,
			CorrelationId: corrId,
			ReplyTo:       q.Name,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	msgs, err := maCh.Consume(q.Name, "recs", true,
		false,
		false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		if corrId == d.CorrelationId {
			var res []int64
			err := json.Unmarshal(d.Body, res)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}

	return nil, nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

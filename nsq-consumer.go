package main

import (
	"encoding/json"
	"log"
	"nsqco/pcerror"

	"github.com/nsqio/go-nsq"
)

type Consumer struct{}

type ErrorMsg struct {
	Error    pcerror.Error `json:"error"`
	UserId   int64         `json:"user_id"`
	GameName string        `json:"game_name"`
}

type NsqMessage struct {
	UserId     int64  `json:"user_id"`
	GameName   string `json:"game_name"`
	Order_Time int64  `json:"order_time"`
}

var url = "pin"
var i int

func (c *Consumer) HandlerForTopicSoloOrder() nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()
		// defer message.Finish()

		var nsqMessage []ErrorMsg
		jsonDecodeErr := json.Unmarshal([]byte(message.Body), &nsqMessage)
		if jsonDecodeErr != nil {
			log.Fatal(jsonDecodeErr)
			message.Finish()
		}

		//模擬API失敗 2 次 第 3 次成功 requeue 5/sec 一次
		// resp := ApiPing(url)
		// if resp.Status != "200 OK" {
		// 	i++
		// 	if i > 1 {
		// 		url = "ping"
		// 	}

		// 	log.Printf("第：%d 次call api 失敗 ， 資料：%s", message.Attempts, nsqMessage)
		// 	message.Requeue(5 * time.Second)
		// }

		log.Println(nsqMessage)
		message.Finish()

		return nil
	})
}

func (c *Consumer) HandlerForTopicOrder() nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()
		defer message.Finish()

		var nsqMessage NsqMessage //[]ErrorMsg //NsqMessage
		jsonDecodeErr := json.Unmarshal([]byte(message.Body), &nsqMessage)
		if jsonDecodeErr != nil {
			log.Fatal(jsonDecodeErr)
		}

		log.Println(nsqMessage)
		return nil
	})
}

func (c *Consumer) HandlerForTopicSageOrder() nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()
		defer message.Finish()

		var nsqMessage []ErrorMsg
		jsonDecodeErr := json.Unmarshal([]byte(message.Body), &nsqMessage)
		if jsonDecodeErr != nil {
			log.Fatal(jsonDecodeErr)
		}

		log.Println(nsqMessage)
		return nil
	})
}

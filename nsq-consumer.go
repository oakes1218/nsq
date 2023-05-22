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

func (c *Consumer) HandlerForTopicSoloOrder() nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()
		defer message.Finish()

		var nsqMessage []ErrorMsg
		jsonDecodeErr := json.Unmarshal([]byte(message.Body), &nsqMessage)
		if jsonDecodeErr != nil {
			log.Fatal(jsonDecodeErr)
		}

		log.Println(nsqMessage, "===")
		return nil
	})
}

func (c *Consumer) HandlerForTopicOrder() nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()
		defer message.Finish()

		var nsqMessage NsqMessage
		jsonDecodeErr := json.Unmarshal([]byte(message.Body), &nsqMessage)
		if jsonDecodeErr != nil {
			log.Fatal(jsonDecodeErr)
		}

		log.Println(nsqMessage, "===")
		return nil
	})
}

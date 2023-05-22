package main

import (
	"errors"
	"fmt"
	"log"
	"nsqco/pclog"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
)

var parser *regexp.Regexp = regexp.MustCompile(`^(\S{3})(.*)`)
var (
	nsqConsumer *NsqConsumer
	nsqProducer *NsqProducer
)

type ParserNSQTopic struct {
	Handler   string `mapstructure:"handler"`
	Topic     string `mapstructure:"topic"`
	Channel   string `mapstructure:"channel"`
	Boardcast string `mapstructure:"boardcast"`
	Ephemeral string `mapstructure:"ephemeral"`
}

type NsqConsumer struct {
	nsqConsumer map[string]nsq.Consumer
}

func (nc *NsqConsumer) SetNsqConsumer(consumer map[string]nsq.Consumer) {
	nc.nsqConsumer = consumer
}

func (nc *NsqConsumer) GettNsqConsumer() map[string]nsq.Consumer {
	return nc.nsqConsumer
}

type NsqProducer struct {
	nsqProducer *nsq.Producer
}

func (np *NsqProducer) SetNsqProducer(producer *nsq.Producer) {
	np.nsqProducer = producer
}

func (np *NsqProducer) GetNsqProducer() *nsq.Producer {
	return np.nsqProducer
}

type Logger interface {
	Output(int, string) error
}

type nsqLog struct {
	Log *pclog.Logger
}

func New(l *pclog.Logger) Logger {
	return &nsqLog{
		Log: l,
	}
}

func (n *nsqLog) Output(i int, s string) (err error) {
	matchs := parser.FindStringSubmatch(s)
	var log, ss string
	if len(matchs) != 3 {
		ss = "ERR"
		log = s
	} else {
		ss = matchs[1]
		log = strings.TrimSpace(matchs[2])
	}

	switch ss {
	case "INF":
		n.Log.Info(log)
	case "WRN":
		n.Log.Warn(log)
	case "ERR":
		n.Log.Error(log)
	default:
		n.Log.Debug(log)
	}
	return
}

func callFuncByName(myClass interface{}, funcName string, params ...interface{}) (ret nsq.Handler, err error) {
	m := reflect.ValueOf(myClass).MethodByName(funcName)
	if !m.IsValid() {
		return nil, errors.New("Method not found " + funcName)
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	out := m.Call(in)

	ret = out[0].Interface().(nsq.Handler)

	return
}

func InitConsumer(tcMap map[string][]string) {
	Nsqlog := New(pclog.Pclog)
	nsqConsumer = &NsqConsumer{}
	consumers := make(map[string]nsq.Consumer)
	for k, v := range tcMap {
		conStruct := &Consumer{}
		fun, err := callFuncByName(conStruct, fmt.Sprintf("HandlerForTopic%s", k))
		if err != nil {
			log.Fatal(err)
		}

		config := nsq.NewConfig()
		config.MaxInFlight = 360
		config.MaxAttempts = 0
		config.MsgTimeout = 10 * time.Minute

		consumer, err := nsq.NewConsumer(v[0], v[1], config)
		if err != nil {
			log.Fatal(err)
		}

		consumer.SetLogger(Nsqlog, nsq.LogLevelInfo)
		consumers[k] = *consumer

		nsqConsumer.SetNsqConsumer(consumers)
		consumer.AddConcurrentHandlers(fun, 200)
		// consumer.AddHandler(soloOrderHandler())

		err = consumer.ConnectToNSQLookupd("nsqlookupd:4161")
		// err = consumer.ConnectToNSQD("127.0.0.1:4150")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func InitProducer() {
	Nsqlog := New(pclog.Pclog)
	nsqProducer = &NsqProducer{}
	config := nsq.NewConfig()
	config.MaxInFlight = 360
	config.MaxAttempts = 0
	config.MsgTimeout = 10 * time.Minute

	producer, err := nsq.NewProducer("test-nsqd:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	producer.SetLogger(Nsqlog, nsq.LogLevelInfo)

	if err := producer.Ping(); err != nil {
		log.Fatal(err)
	}

	nsqProducer.SetNsqProducer(producer)
}

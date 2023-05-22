package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nsqco/pcerror"
	"nsqco/pclog"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const InErrService = pcerror.TEST

var Val *Config

type Config struct {
	ServerPort      string                    `mapstructure:"SERVER_PORT" json:"SERVER_PORT"`
	NsqReceiveTopic map[string]ParserNSQTopic `mapstructure:"NSQTOPIC"`
}

func init() {
	//讀取還近變數
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PC")

	//取設定黨
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := viper.Unmarshal(&Val)
	if err != nil {
		panic(err)
	}

	// driver, ok := os.LookupEnv("TERM")

	log.Println("ENV:", Val)
	log.Println("Cofing 設定成功")
	log.Println(viper.Get("NSQTOPIC"))
	log.Println(viper.Get("NSQRECEIVETOPIC_TOPIC_2_HANDLER"), "")
	// init log
	pclog.Pclog = pclog.New()
	// init user DB conn
	// model.InitUser()
}

func main() {
	// pclog.Pclog.WithError(pcerror.Wrap(1000, "test", InErrService, pcerror.Extrainfo, errors.New("test"))).Info("qwe")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		close(sigs)
		for _, v := range nsqConsumer.GettNsqConsumer() {
			v.Stop()
		}
		nsqProducer.GetNsqProducer().Stop()

		i := 1
		for {
			time.Sleep(time.Second * 1)
			if i == 1 {
				break
			}
			i--
		}

		done <- true
		close(done)
	}()

	InitNsq()
	httpServer()

	log.Println("exit")
}

func InitNsq() {
	tcMap := make(map[string][]string)
	for _, v := range Val.NsqReceiveTopic {
		tcS := []string{v.Topic, v.Channel}
		tcMap[v.Handler] = tcS
	}

	InitConsumer(tcMap)
	InitProducer()

	time.Sleep(3 * time.Second)
	messageBody := []byte(`[{"error":{"msg":"SYSTEM_ERROR","code":125041045,"extrainfo":null,"time":1645692452,"service":"SOLO","origin_err":"Key"},"user_id":1095926084804800513,"game_name":""}]`)
	topicName := "SOLO_ORDER_ERROR_TOPIC"
	err := nsqProducer.GetNsqProducer().Publish(topicName, messageBody)
	if err != nil {
		log.Fatal(err)
	}

	messageBody1 := []byte(` {"order_id":1036904237549768705,"order_mode":1,"game_manager_id":1,"platform_id":10001,"user_id":12345678901,"game_manager_name":"GAMETE01","platform_name":"MER10001","user_name":"MER10001@test1","game_name":"BBQL","draw_id":1031,"draw_num":"201809041031","rule":"r0013","tag":"","choose":"LOCATE:S:ODD","extra_info":null,"total_bet":1,"total_bet_gold":5,"odds":{"LOCATE:S:ODD":1.94},"fair_odds":{"LOCATE:S:ODD":2},"origin_odds":{"LOCATE:S:ODD":1.94},"bet_gold":5,"win_gold":0,"paid_gold":0,"profit_gold":0,"wars":0,"pay_max":500000,"result":null,"result_display":null,"status":1,"currency":"CNY","exchange_rate":1,"win":null,"lose":null,"tie":null,"error_code":0,"ip":"202.3.3.1","entrance":1,"portal":1,"client":1,"device":1,"order_time":1536052255,"computed_time":0,"computed_count":0,"updated_time":1536052209,"cart_id":1036904237549768704,"trace_id":"thisIsTestingTraceIDCancel","is_auto":false}`)
	topicName1 := "EDEN_ORDER_CREDIT"
	err = nsqProducer.GetNsqProducer().Publish(topicName1, messageBody1)
	if err != nil {
		log.Fatal(err)
	}
}

func httpServer() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", Pong)
	r.GET("/readiness", Pong)

	r.Run(string(fmt.Sprintf("%v", viper.Get("SERVER_PORT"))))
}

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gopkg.in/Shopify/sarama.v1"

	_dbpollingService "github.com/gamepkw/db-polling-microservice/internal/services"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// logger.Info("start program...")

	config := sarama.NewConfig()
	config.ClientID = "my-kafka-client"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	brokers := []string{viper.GetString("kafka.broker_address")}
	kafkaClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaClient.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001", "http://localhost:3000", "http://localhost:8090"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// go tu.ConsumeScheduledTransaction(ctx)

	pollingInterval := 15 * time.Second

	pu := _dbpollingService.NewPollingService(pollingInterval)
	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go pu.Polling(ctx, &wg, pollingInterval, stopChan)

	log.Fatal(e.Start(viper.GetString("server.address")))

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
	time.Sleep(1 * time.Minute)

	close(stopChan)

	wg.Wait()

}

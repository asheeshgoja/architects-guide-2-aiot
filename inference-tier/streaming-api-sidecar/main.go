package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"os"
	"strconv"
	"time"
	// "bufio"
	"log"
)


func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok { 
		return defaultValue
	}
	fmt.Println("Env var ",envVar , " = ",envVarValue)
	return envVarValue
}

func lookupInt64Env(envVar string, defaultValue int64) int64 {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	int64Val, _ := strconv.ParseInt(envVarValue, 10, 64)
	return int64Val
}

const (

	connType      = "tcp"
	stopCharacter = "\r\n\r\n"
)


func main() {

	startKafkaConsumer()

}

func startKafkaConsumer() {

	connHost   := "127.0.0.1"
	connPort   := lookupStringEnv("PORT", "9898")

	topic         := lookupStringEnv("TOPIC", "shaded-pole-motor-sensor_data")
	brokerAddress := lookupStringEnv("KAFKA-BROKER", "35.236.22.237:32199")


	ctx := context.Background()

	// l := log.New(os.Stdout, "kafka reader: ", 0)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// GroupID: lookupStringEnv("STREAM_GRP_ID","work-ld-1"),
		// assign the logger to the reader
		// Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		fmt.Printf("Waiting for kafka mesages from broker %s on topic %s\n" , brokerAddress, topic)

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))


		retry(100, 4, func() error {
			var socketConnection net.Conn

			fmt.Printf("starting sidecar socket on port %s \n", connPort)
			socketConnection, err := net.Dial(connType, connHost+":"+connPort)

			if err != nil {
				fmt.Println("Error connecting:, retrying ....", err.Error())
				return err
			}

			buff := []byte(msg.Value)
			_, e := socketConnection.Write([]byte(buff))
			// _ , e = socketConnection.Write([]byte(stopCharacter))

			if e != nil {
				log.Print("could not send message " + e.Error())
			}

			// message, _ := bufio.NewReader(socketConnection).ReadString('\n')
			// log.Print("socket Server relay: " + message)

			buff2 := make([]byte, 1024)
			n, _ := socketConnection.Read(buff2)
			log.Printf("Receive: %s", buff2[:n])

			socketConnection.Close()
			return nil
		})
	}

}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep * time.Second)
			return retry(attempts, 1*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

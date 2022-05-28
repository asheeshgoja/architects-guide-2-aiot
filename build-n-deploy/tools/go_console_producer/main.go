package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

/*
Options:
 [-b]                 broker
 [-t]                 topic
*/

// default values for environment variables
const (
	BootstrapServersDefault = "35.236.22.237:32199"
	TopicDefault            = "control-message"
	DelayDefault            = 1000
	MessageDefault          = ""
	MessageCountDefault     = 1000000
	ProducerAcksDefault     = int16(1)
)

// ProducerConfig defines the producer configuration
type ProducerConfig struct {
	BootstrapServers string
	Topic            string
	Delay            int
	Message          string
	MessageCount     int64
	ProducerAcks     int16
}

func NewProducerConfig() *ProducerConfig {
	topic := flag.String("t", "shaded-pole-motor-sensor_data", "")
	broker := flag.String("b", "35.236.22.237:32199", "")
	flag.Parse()

	config := ProducerConfig{
		BootstrapServers: *broker,
		Topic:            *topic,
		Delay:            DelayDefault,
		Message:          MessageDefault,
		MessageCount:     MessageCountDefault,
		ProducerAcks:     ProducerAcksDefault,
	}
	return &config
}

func main() {

	config := NewProducerConfig()
	log.Printf("Go producer starting with config=%+v\n", config)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
	if err != nil {
		log.Printf("Error creating the Sarama sync producer: %v", err)
		os.Exit(1)
	}

	end := make(chan int, 1)
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			produceMessage(scanner.Text(), config.Topic, producer)
		}
	}()

	// waiting for the end of all messages sent or an OS signal
	select {
	case <-end:
		log.Printf("Finished to send %d messages\n", config.MessageCount)
	case sig := <-signals:
		log.Printf("Got signal: %v\n", sig)
	}

	err = producer.Close()
	if err != nil {
		log.Printf("Error closing the Sarama sync producer: %v", err)
		os.Exit(1)
	}
	log.Printf("Producer closed")
}

func produceMessage(messageVal string, topic string, producer sarama.SyncProducer) {
	// ipAdd := GetLocalIP()
	// value := fmt.Sprintf(messageVal)
	log.Printf("Captured message from console : ", messageVal)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageVal),
	}
	log.Printf("Sending message: value=%s\n", messageVal)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Erros sending message: %v\n", err)
	} else {
		log.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
	}

}

//mosquitto_pub -h "35.236.22.237" -t shaded-pole-motor-sensor_data -m "{"deviceID": "14333616", "current": 4.69, "temperature": 39.06, "vibration": 39.06, "sound": 12.50, "fft_data": "true"}"
//mosquitto_pub -h "35.236.22.237" -t control-message -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
//mosquitto_sub -h "35.236.22.237" -t control-message
//mosquitto_sub -h "35.236.22.237" -t shaded-pole-motor-sensor_data
// go_console_producer# go run . -t control-message {"command": "download-model", "payload" : "2022-05-06-20:36:11-model.tflite"}
// go_console_producer# go run . -t control-message {"command": "train-model", "payload" : "normalized_training_data_2022-05-16T14:55:39:1652712939683.csv"}
// go_console_producer# go run . -t control-message {"command": "extract-data", "payload" : "raw_sensor_training_data_2022-05-26T23:15:39:1653606939487.csv"}

//hivemq
// mqtt pub -h 35.236.22.237 -p 30005 -t control-message  -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
// mqtt pub -h 35.236.22.237 -p 30005 -t shaded-pole-motor-sensor_data -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
// mqtt sub -h 35.236.22.237 -p 30005 -t shaded-pole-motor-sensor_data

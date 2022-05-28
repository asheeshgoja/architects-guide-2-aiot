package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/segmentio/kafka-go"
)

func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	fmt.Println("Env var ", envVar, " = ", envVarValue)
	return envVarValue
}

func main() {

	SetupCloseHandler()

	mqtt_subs_opts := MQTT.NewClientOptions()
	initMqttBroker(mqtt_subs_opts)
	go initKafkaBroker(mqtt_subs_opts)

	startKafkaConsumer()

}

func initMqttBroker(opts *MQTT.ClientOptions) {
	mqtt_broker := lookupStringEnv("MQTT-BROKER", "tcp://35.236.22.237:1883")
	// mqtt_password := lookupStringEnv("MQTT-USER", "")
	// mqtt_user := lookupStringEnv("MQTT-PSW", "")
	mqtt_id := lookupStringEnv("MQTT-ID", "architectsguide2aiot_mqtt-id")
	cleansess := false
	store := ":memory:"

	opts.AddBroker(mqtt_broker)
	opts.SetClientID(mqtt_id)
	// opts.SetUsername(mqtt_user)
	// opts.SetPassword(mqtt_password)
	opts.SetCleanSession(cleansess)
	if store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(store))
	}

	fmt.Printf("Initilized mqtt client at broker %s \n", mqtt_broker)
}

func initKafkaBroker(mqtt_opts *MQTT.ClientOptions) {

	kafka_topic := lookupStringEnv("DATA-TOPIC", "shaded-pole-motor-sensor_data")
	kafka_brokerAddress := lookupStringEnv("KAFKA-BROKER", "35.236.22.237:32199")

	// config := NewProducerConfig()
	// fmt.Printf("Go producer starting with config=%+v\n", config)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(int16(1))
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafka_brokerAddress}, producerConfig)
	if err != nil {
		fmt.Printf("Error creating the Sarama sync producer: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Started kafka producer client at broker %s , topic %s \n", kafka_brokerAddress, kafka_topic)

	end := make(chan int, 1)
	// go func() {

	// 	startMqttSubscriber(mqtt_opts, )
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	for scanner.Scan() {
	// 		publishMqttMessageToKafka(scanner.Text(), kafka_topic, producer)
	// 	}
	// }()

	startMqttSubscriber(mqtt_opts, func(messageVal string) {
		// fmt.Printf("Captured message from console : ", messageVal)

		msg := &sarama.ProducerMessage{
			Topic: kafka_topic,
			Value: sarama.StringEncoder(messageVal),
		}
		fmt.Printf("Sending message: value=%s\n", messageVal)
		// partition, offset, err := producer.SendMessage(msg)
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Printf("Erros sending message: %v\n", err)
		}
		// else {
		// 	fmt.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
		// }
	})

	// waiting for the end of all messages sent or an OS signal
	select {
	case <-end:
		fmt.Printf("Finished to send  messages\n")
	case sig := <-signals:
		fmt.Printf("Got signal: %v\n", sig)
	}

	err = producer.Close()
	if err != nil {
		fmt.Printf("Error closing the Sarama sync producer: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Producer closed")

}

func startMqttSubscriber(opts *MQTT.ClientOptions, publishMqttMessageToKafka func(string)) {
	qos := 0
	mqtt_topic := lookupStringEnv("DATA-TOPIC", "shaded-pole-motor-sensor_data")

	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(mqtt_topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		publishMqttMessageToKafka(incoming[1])
	}

	client.Disconnect(250)
	fmt.Println("mqtt-kafka-protocol-bridge disconnected")

}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func startKafkaConsumer() {

	mqtt_broker := lookupStringEnv("MQTT-BROKER", "tcp://35.236.22.237:1883")
	kafka_broker := lookupStringEnv("KAFKA-BROKER", "35.236.22.237:32199")

	kafka_topic := lookupStringEnv("CONTROL-TOPIC", "control-message")
	mqtt_pub_topic := lookupStringEnv("CONTROL-TOPIC", "control-message")
	qos := 0

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqtt_broker)
	opts.SetClientID("test")
	// opts.SetUsername(*user)
	// opts.SetPassword(*password)
	opts.SetCleanSession(false)

	mqtt_client := MQTT.NewClient(opts)
	if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	ctx := context.Background()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafka_broker},
		Topic:   kafka_topic,
		// GroupID: lookupStringEnv("STREAM_GRP_ID","work-ld-1"),
		// assign the logger to the reader
		// Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		fmt.Printf("Waiting for kafka mesages from broker %s on topic %s\n", kafka_broker, kafka_topic)

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

		retry(100, 4, func() error {

			fmt.Println("Sample Publisher Started")
			token := mqtt_client.Publish(mqtt_pub_topic, byte(qos), false, string(msg.Value))
			token.Wait()
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

//mosquitto_pub -h "35.236.22.237" -t shaded-pole-motor-sensor_data -m "{"deviceID": "14333616", "current": 4.69, "temperature": 39.06, "vibration": 39.06, "sound": 12.50, "fft_data": "true"}"
//mosquitto_pub -h "35.236.22.237" -t control-message -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
//mosquitto_sub -h "35.236.22.237" -t control-message
//mosquitto_sub -h "35.236.22.237" -t shaded-pole-motor-sensor_data
// go_console_producer# go run . -t control-message {"command": "download-model", "payload" : "2022-05-06-20:36:11-model.tflite"}
// go_console_producer# go run . -t control-message {"command": "train-model", "payload" : "normalized_training_data_2022-05-16T14:55:39:1652712939683.csv"}
// go_console_producer# go run . -t control-message {"command": "extract-data", "payload" : "raw_sensor_training_data_2022-05-06T16:14:48:1651853688477.csv"}

//hivemq
// mqtt pub -h 35.236.22.237 -p 30005 -t control-message  -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
// mqtt pub -h 35.236.22.237 -p 30005 -t shaded-pole-motor-sensor_data -m "{\"command\": \"download-model\", \"payload\" : \"2022-04-25-18:39:17-model.tflite\"}"
// mqtt sub -h 35.236.22.237 -p 30005 -t shaded-pole-motor-sensor_data

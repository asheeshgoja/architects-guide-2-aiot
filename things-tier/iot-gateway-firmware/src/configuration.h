//Wifi credentials
const char *ssid = "";
const char *pass = "";


//MQTT Broker
const char *broker = "10.0.0.30";
const int  broker_port = 30005;
const char *brokerUser = "YOURUSER";
const char *brokerPass = "YOURPASSWORD";

//Data message topic
const char *dataTopic = "shaded-pole-motor-sensor_data";

//Control Topic
const char *commandTopic = "control-message";

//Model OTA URL
const char *modelRegistry = "https://10.0.0.30:30007/quantized/";

// Device Activation URL
const char *deviceRegistry = "https://10.0.0.30:30006/confirmActivation";
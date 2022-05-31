#include <Arduino.h>
#include <WiFi.h>
#include <WiFiMulti.h>
#include "fft_module.h"
#include <ArduinoJson.h>
#include <HTTPClient.h>

#include "tflm_module.h"
#include "motor_sensors.h"
#include "mqtt_module.h"
#include "fft_module.h"
#include "servo_module.h"
#include "aggregation.h"
#include "activation_module.h"
// #include "configuration_local.h"

#include "configuration.h"

const int TMP_PIN = 36;
const int VIBRATION_PIN = 39;
const int SOUND_PIN = 34;
const int CURRENT_PIN = 35;
const int LED_BUILTIN = 2;
const int SERVO_PIN = 4;

HTTPClient http_client_activation, http_client_model;
WiFiClient espClient;
long currentTime, lastTime;
int isValidDevice = 0;

TFLM_Module tflm_module;
Sensors_Module sensors_module;
Mqtt_Module mqtt_module;
Fft_Module fft_module;
Servo_Module hydraulicValveController(SERVO_PIN);
Aggregation_Module aggregationModule;
Activation_Module activationModule;

uint32_t chipId = 0;

void setupWifi()
{
  delay(100);
  Serial.print("\nConnecting to");
  Serial.println(ssid);

  WiFi.begin(ssid, pass);

  while (WiFi.status() != WL_CONNECTED)
  {
    delay(100);
    Serial.print("-");
  }

  Serial.print("\nConnected to ");
  Serial.println(ssid);
}

void printBinaryData(byte *buf, int len)
{
  Serial.print("[");
  for (int x = 0; x < len; x++)
  {
    // Serial.println(buf[x]);
    Serial.print(buf[x], BIN);
    // Serial.print("\t");
  }
  Serial.println("]");
}

void setup()
{
  Serial.begin(115200);

  pinMode(LED_BUILTIN, OUTPUT);

  setupWifi();

  activationModule.init(&http_client_activation, deviceRegistry);
  chipId = activationModule.getChipID();
  isValidDevice = activationModule.isDeviceActivated();
  Serial.printf("Device Activation code from server = %d\n", isValidDevice);
  if (0 == isValidDevice)
  {
    Serial.printf("Device Activation failed!!\n");
  }

  tflm_module.init(&http_client_model, modelRegistry);
  sensors_module.init(TMP_PIN, CURRENT_PIN, VIBRATION_PIN, SOUND_PIN);
  mqtt_module.init(espClient, broker, broker_port, commandTopic, dataTopic, String(chipId));

  hydraulicValveController.turnValveOff();
}

void loop()
{
  if (isValidDevice == 0)
  {
    Serial.printf("Device not registered, retrying ... !\n");
    delay(3000);
    isValidDevice = activationModule.isDeviceActivated();
    return;
  }

  mqtt_module.reconnect();

  float tempVal, vibrationVal, soundVal, currentVal;
  tempVal = vibrationVal = soundVal = currentVal = 0;

  int maxSampleSize = fft_module.getMaxSampleSize();
  int samplingPeriod = fft_module.getSamplingPeriod();

  long sampleStartMicros = 0;

  for (int i = 0; i < maxSampleSize; i++)
  {
    sampleStartMicros = micros();
    aggregationModule.aggregateData(sensors_module, currentVal, tempVal, vibrationVal, soundVal);

    char sensor_data[255];
    snprintf(sensor_data, 255, "{\"deviceID\": \"%d\", \"current\": %.2f, \"temperature\": %.2f, \"vibration\": %.2f, \"sound\": %.2f}",
             chipId, currentVal, tempVal, vibrationVal, soundVal);

    Serial.println(sensor_data);
    // mqtt_module.publish(sensor_data);

    fft_module.recordSample(currentVal, tempVal, vibrationVal, soundVal);
    while ((micros() - sampleStartMicros) < samplingPeriod)
    { /* spin */
    }
  }

  float temp_fft, current_fft, vibr_fft, sound_fft;
  temp_fft = current_fft = vibr_fft = sound_fft = 0;

  fft_module.perform_fft(&current_fft, &temp_fft, &vibr_fft, &sound_fft);

  char sensor_data_fft[255];
  snprintf(sensor_data_fft, 255, "{\"deviceID\": \"%d\", \"current\": %.2f, \"temperature\": %.2f, \"vibration\": %.2f, \"sound\": %.2f, \"fft_data\": \"true\"}",
           chipId, current_fft, temp_fft, vibr_fft, sound_fft);

  Serial.println(sensor_data_fft);
  mqtt_module.publish(sensor_data_fft);

  fft_module.resetSampleCounter();

  float inference_val = tflm_module.predict(current_fft, temp_fft, vibr_fft, sound_fft);

  // Serial.printf("performing tflm inference on the  result = %.2f \n", result);

  if (inference_val > 8.0)
  {
    hydraulicValveController.turnValveOn();
  }
}
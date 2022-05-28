#include "mqtt_module.h"
#include <Arduino.h>
#include "tflm_module.h"
#include <ArduinoJson.h>

extern uint8_t *g_tflileModelBuf;

void Mqtt_Module::init(Client &client, const char *broker, int port, String sub_topic, String pub_topic, String id)
{
    ptr_pub_sub_client = new PubSubClient(client);

    ptr_pub_sub_client->setServer(broker, port);
    ptr_pub_sub_client->setCallback(Mqtt_Module::callback);

    mqtt_broker = broker;
    subscription_topic = sub_topic;
    publication_topic = pub_topic;
    client_id = id;
}

void Mqtt_Module::reconnect()
{
    if (!ptr_pub_sub_client->connected())
    {
        while (!ptr_pub_sub_client->connected())
        {
            Serial.printf("Connecting to broker %s , clientid : %s ,  on topic %s\n", mqtt_broker.c_str(), client_id.c_str(),  subscription_topic.c_str());

            // if (ptr_pub_sub_client->connect("architectsguide2aiot-aiot-demo", brokerUser, brokerPass))
            if (ptr_pub_sub_client->connect(client_id.c_str()))
            {
                Serial.printf("Connecting to broker %s , clientid : %s ,  on topic %s\n", mqtt_broker.c_str(), client_id.c_str(),  subscription_topic.c_str());
                ptr_pub_sub_client->subscribe(subscription_topic.c_str());
            }
            else
            {
                Serial.printf("Trying to connect to broker %s , clientid : %s ,  on topic %s\n", mqtt_broker.c_str(), client_id.c_str(),  subscription_topic.c_str());
                delay(5000);
            }
        }
    }

    ptr_pub_sub_client->loop();
}

void Mqtt_Module::callback(char *topic, uint8_t *payload, unsigned int length)
{

    char buf[255] = "";
    for (int i = 0; i < length; i++)
    {
        // Serial.print((char)payload[i]);
        buf[i] = (char)payload[i];
    }

    Serial.printf("mqtt message received from broker on topic %s , payload : %s \n", topic, buf);

    String jsonMessage(buf);
    // Serial.println(jsonMessage);

    //   Deserialize the JSON document
    StaticJsonDocument<255> jsonBuffer;
    DeserializationError error = deserializeJson(jsonBuffer, jsonMessage);

    // Test if parsing succeeds.
    if (error)
    {
        Serial.printf("deserializeJson() failed: %s \n", error.c_str());
        return;
    }

    const char *command = jsonBuffer["command"];
    const char *cmd_payload = jsonBuffer["payload"];

    Serial.printf("command:%s , payload: %s\n", command, cmd_payload);

    if (strcmp(command, "download-model") == 0)
    {
        for (int i = 0; i < 30; i++)
        {
            digitalWrite(2, HIGH);
            delay(100);
            digitalWrite(2, LOW);
            delay(100);
        }

        Serial.println();
        TFLM_Module::setNewModelFileName(cmd_payload);
    }
}

void Mqtt_Module::publish(char *payload)
{
    boolean r = ptr_pub_sub_client->publish(publication_topic.c_str(), payload);
    r ? Serial.println("Publishing data success ") : Serial.println("Publishing data failed!");
}
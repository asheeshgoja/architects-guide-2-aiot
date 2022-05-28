#ifndef ___MQTT_MODULE__
#define ___MQTT_MODULE__

#include <stdint.h>
#include <PubSubClient.h>

class PubSubClient;
class Client;

class Mqtt_Module
{

private:
    PubSubClient* ptr_pub_sub_client;
    String mqtt_broker;
    String subscription_topic;
    String publication_topic;
    String client_id;

public:
    void init(Client& client,  const char*  broker,int port, String sub_topic, String pub_topic, String id );
    void reconnect();
    void publish(char *payload);

private:
    static void callback(char *topic, uint8_t *payload, unsigned int length);
};

#endif
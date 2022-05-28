#include "activation_module.h"
#include <HTTPClient.h>

extern const char* reg_svr_pub_key;

void Activation_Module::init(HTTPClient *p_http, String svr)
{
    activationServer = svr;
    pHTTPClient = p_http;
}

uint32_t Activation_Module::getChipID()
{
  uint32_t c_id = 0;

  for (int i = 0; i < 17; i = i + 8)
  {
    c_id |= ((ESP.getEfuseMac() >> (40 - i)) & 0xff) << i;
  }

  Serial.printf("Device ID : %d\n", c_id);
  return c_id;
}

int Activation_Module::isDeviceActivated()
{
    uint32_t chipID = getChipID();

    Serial.printf("Connecting to device registry server @%s\n", activationServer.c_str());
    pHTTPClient->begin(activationServer.c_str(), reg_svr_pub_key);

    pHTTPClient->addHeader("Content-Type", "application/x-www-form-urlencoded");

    char payload[255] = {};
    snprintf(payload, 255, "%s=%d", "device_id", chipID);

    Serial.printf("posting payload %s\n", payload);
    int httpResponseCode = pHTTPClient->POST(payload);

    if (httpResponseCode == 200)
    {
        Serial.println(httpResponseCode);
        String activationCode = pHTTPClient->getString();
        Serial.printf("Activation code %s\n", activationCode.c_str());
        pHTTPClient->end();
        return strcmp(activationCode.c_str() , "TRUE") == 0 ? 1 : 0;
    }
    else
    {
        Serial.printf("Activation failed!!\n");
        Serial.println(httpResponseCode);
        pHTTPClient->end();
        return 0;
    }

    return 0;
}
#ifndef __ACTIVATION_MODULE__
#define __ACTIVATION_MODULE__

#include <Arduino.h>

class HTTPClient;

class Activation_Module
{
    HTTPClient *pHTTPClient;
    String activationServer;

public:
    int isDeviceActivated();
    void init(HTTPClient *p_http, String svr);
    uint32_t getChipID();
};

#endif
#ifndef ___SENSORS__
#define ___SENSORS__

class Sensors_Module
{
    int temperaturePin;
    int soundPin;
    int vibrationPin;
    int currentPin;

public:
    void init(int temp_pin, int current_pin, int vibr_pin, int sound_pin);
    double getCurrent();
    double  getSound();
    double  getVibration();
    double  getTemperature();
};


#endif

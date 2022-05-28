#include "aggregation.h"
#include "motor_sensors.h"

void Aggregation_Module::aggregateData(Sensors_Module &sensorModule, float &current, float &temperature, float &vibration, float &sound)
{
     temperature = sensorModule.getTemperature();
     vibration = sensorModule.getVibration();
     sound = sensorModule.getSound();
     current = sensorModule.getCurrent();
}
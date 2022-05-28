#include "servo_module.h"
#include <Servo.h>

Servo_Module::Servo_Module(int pin)
{
    pHydraulic_valve_servo = new Servo();
    pHydraulic_valve_servo->attach(pin);
}

void Servo_Module::turnValveOn()
{
    for (int posDegrees = 0; posDegrees <= 180; posDegrees++)
    {
        pHydraulic_valve_servo->write(posDegrees);
        Serial.println(posDegrees);
        delay(20);
    }
    Serial.printf("hydraulic valve servo turned on !\n");
}

void Servo_Module::turnValveOff()
{
    for (int posDegrees = 180; posDegrees >= 0; posDegrees--)
    {
        pHydraulic_valve_servo->write(posDegrees);
        Serial.println(posDegrees);
        delay(20);
    }
    Serial.printf("hydraulic valve servo turned off !\n");
}
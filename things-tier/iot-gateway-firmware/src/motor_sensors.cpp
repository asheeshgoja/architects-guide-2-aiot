#include "motor_sensors.h"
#include <Arduino.h>
#include "EmonLib.h"

EnergyMonitor emon1;

void Sensors_Module::init(int temp_pin, int current_pin, int vibr_pin, int sound_pin)
{
    temperaturePin = temp_pin;
    soundPin = current_pin;
    vibrationPin = vibr_pin;
    currentPin = sound_pin;

    emon1.current(current_pin, 111.1);
}

double Sensors_Module::getCurrent()
{
    double Irms = emon1.calcIrms(1480); // Calculate Irms only
                                        // emon1.calcVI(20, 2000);

    // Serial.print(currentVal * 230.0); // Apparent power (230 is the Voltage)
    // Serial.print(" ");
    // Serial.println(currentVal); // Irms

    float realPower = emon1.realPower;         //extract Real Power into variable
    float apparentPower = emon1.apparentPower; //extract Apparent Power into variable
    float powerFActor = emon1.powerFactor;     //extract Power Factor into Variable
    float supplyVoltage = emon1.Vrms;          //extract Vrms into Variable
    float currentVal = Irms;

    // Serial.printf("realPower: %.2f, apparentPower: %.2f, powerFActor: %.2f, supplyVoltage: %.2f, Irms: %.2f\n",
    //                       realPower, apparentPower, powerFActor, supplyVoltage, Irms);

    return (Irms * 100) + 400; // add bias V
}

double Sensors_Module::getSound()
{
    double sensorValue = analogRead(soundPin);

    unsigned int peakToPeak = 0; // peak-to-peak level

    unsigned int signalMax = 0;
    unsigned int signalMin = 1024;

    float sample = sensorValue;
    if (sample < 1024) // toss out spurious readings
    {
        if (sample > signalMax)
        {
            signalMax = sample; // save just the max levels
        }
        else if (sample < signalMin)
        {
            signalMin = sample; // save just the min levels
        }
    }

    peakToPeak = signalMax - signalMin;      // max - min = peak-peak amplitude
    float volts = (peakToPeak * 5.0) / 1024; // convert to volts

    // Serial.printf("Raw sound sensor val %.2f , normalizedVAl %.2f\n", sensorValue, volts);

    return volts;
}

double Sensors_Module::getVibration()
{
    return analogRead(vibrationPin);
}

double Sensors_Module::getTemperature()
{
    float tempVal =  analogRead(temperaturePin);
    float volts = tempVal / 1023.0;
    float temp = (volts - 0.5) * 100;
    return temp;
}
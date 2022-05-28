#ifndef ___SERVO_MODULE__
#define ___SERVO_MODULE__

class Servo;

class Servo_Module
{
    Servo* pHydraulic_valve_servo; 
public:
    Servo_Module(int pin);
    void turnValveOn();
    void turnValveOff();
};


#endif
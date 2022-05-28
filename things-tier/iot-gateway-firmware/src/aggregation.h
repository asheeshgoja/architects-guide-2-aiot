#ifndef __AGGRE_MODULE__
#define __AGGRE_MODULE__

class Sensors_Module;

class Aggregation_Module
{

public:
    void aggregateData(Sensors_Module &sensorModule, float& current, float& temperature, float& vibration, float& sound );

};


#endif
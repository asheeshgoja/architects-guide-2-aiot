#ifndef __ARDUINO_FFT_MODULE__
#define __ARDUINO_FFT_MODULE__

#include "arduinoFFT.h"

const int MAX_SAMPLE_SIZE = 64;

class  arduinoFFT;

class Fft_Module
{

private:
    // const uint16_t samples = 64; // This value MUST ALWAYS be a power of 2
    const double signalFrequency = 1000;
    const double samplingFrequency = 13095; // find using the function getMaxSamplingFrequency
    const uint8_t amplitude = 100;

    // double* fft_input_data_temperature;
    // double* fft_input_data_current;
    // double* fft_input_data_vibration;
    // double* fft_input_data_sound;
    // double* vImag;

    double fft_input_data_temperature[MAX_SAMPLE_SIZE] = {};
    double fft_input_data_current[MAX_SAMPLE_SIZE] = {};
    double fft_input_data_vibration[MAX_SAMPLE_SIZE] = {};
    double fft_input_data_sound[MAX_SAMPLE_SIZE] = {};
    double vImag[MAX_SAMPLE_SIZE] = {};


    int sample_counter = 0;

    arduinoFFT* ptrFFT;
    // arduinoFFT* ptrFFT_current;
    // arduinoFFT* ptrFFT_vibration;
    // arduinoFFT* ptrFFT_sound;

public:
    Fft_Module();
    void perform_fft(float *current_fft, float *temperature_fft, float *vibration_fft, float *sound_fft);
    void recordSample(float current, float temperature, float vibration, float sound);
    int getSampleCounterVal();
    int getMaxSampleSize();
    void resetSampleCounter();
    double doFftOnOneScalarAndGetMajorPeak(arduinoFFT *ptrFFT, double* vReal);
    double getSamplingPeriod();
    double majorPeakParabola(double *vReal);
  	void Parabola(double x1, double y1, double x2, double y2, double x3, double y3, double *a, double *b, double *c);

};

#endif
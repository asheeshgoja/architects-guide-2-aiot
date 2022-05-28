#include "fft_module.h"

#define SCL_INDEX 0x00
#define SCL_TIME 0x01
#define SCL_FREQUENCY 0x02
#define SCL_PLOT 0x03

Fft_Module::Fft_Module()
{
    // ptrFFT = new arduinoFFT(fft_input_data_current, vImag, MAX_SAMPLE_SIZE, samplingFrequency);
    ptrFFT = new arduinoFFT();
}

void Fft_Module::perform_fft(float *current_fft, float *temperature_fft, float *vibration_fft, float *sound_fft)
{
    *current_fft = doFftOnOneScalarAndGetMajorPeak(ptrFFT, fft_input_data_current);
    memset(vImag, 0, sizeof(vImag));

    *temperature_fft = doFftOnOneScalarAndGetMajorPeak(ptrFFT, fft_input_data_temperature);
    memset(vImag, 0, sizeof(vImag));

    *vibration_fft = doFftOnOneScalarAndGetMajorPeak(ptrFFT, fft_input_data_vibration);
    memset(vImag, 0, sizeof(vImag));

    *sound_fft = doFftOnOneScalarAndGetMajorPeak(ptrFFT, fft_input_data_sound);
}

double Fft_Module::getSamplingPeriod()
{
    return round(1000000 * (1.0 / samplingFrequency));
}

double Fft_Module::doFftOnOneScalarAndGetMajorPeak(arduinoFFT *ptrFFT, double *vReal)
{
    ptrFFT->Windowing(vReal, MAX_SAMPLE_SIZE, FFT_WIN_TYP_HAMMING, FFT_FORWARD); /* Weigh data */
    ptrFFT->Compute(vReal, vImag, MAX_SAMPLE_SIZE, FFT_FORWARD);                 /* Compute FFT */
    ptrFFT->ComplexToMagnitude(vReal, vImag, MAX_SAMPLE_SIZE);                   /* Compute magnitudes */

    double peak = majorPeakParabola(vReal);
    return peak;
}

void Fft_Module::recordSample(float current, float temperature, float vibration, float sound)
{
    fft_input_data_sound[sample_counter] = sound;
    fft_input_data_temperature[sample_counter] = temperature;
    fft_input_data_current[sample_counter] = current;
    fft_input_data_vibration[sample_counter] = vibration;

    vImag[sample_counter] = 0.0;
    sample_counter++;
}

int Fft_Module::getSampleCounterVal()
{

    return sample_counter;
}

void Fft_Module::resetSampleCounter()
{
    sample_counter = 0;
}

int Fft_Module::getMaxSampleSize()
{
    return MAX_SAMPLE_SIZE;
}

int getMaxSamplingFrequency()
{
    long newTime = micros();
    for (int i = 0; i < 1000000; i++)
    {
        analogRead(34);
    }
    float t = (micros() - newTime) / 1000000.0;

    return (1.0 / t) * 1000000;
}

double Fft_Module::majorPeakParabola(double *vReal)
{
    double maxY = 0;
    uint16_t IndexOfMaxY = 0;
    // If sampling_frequency = 2 * max_frequency in signal,
    // value would be stored at position samples/2
    for (uint16_t i = 1; i < ((MAX_SAMPLE_SIZE >> 1) + 1); i++)
    {
        if ((vReal[i - 1] < vReal[i]) && (vReal[i] > vReal[i + 1]))
        {
            if (vReal[i] > maxY)
            {
                maxY = vReal[i];
                IndexOfMaxY = i;
            }
        }
    }

    double freq = 0;
    if (IndexOfMaxY > 0)
    {
        // Assume the three points to be on a parabola
        double a, b, c;
        Parabola(IndexOfMaxY - 1, vReal[IndexOfMaxY - 1], IndexOfMaxY, vReal[IndexOfMaxY], IndexOfMaxY + 1, vReal[IndexOfMaxY + 1], &a, &b, &c);

        // Peak is at the middle of the parabola
        double x = -b / (2 * a);

        // And magnitude is at the extrema of the parabola if you want It...
        // double y = a*x*x+b*x+c;

        // Convert to frequency
        freq = (x * samplingFrequency) / (MAX_SAMPLE_SIZE);
    }

    return freq;
}

void Fft_Module::Parabola(double x1, double y1, double x2, double y2, double x3, double y3, double *a, double *b, double *c)
{
    double reversed_denom = 1 / ((x1 - x2) * (x1 - x3) * (x2 - x3));

    *a = (x3 * (y2 - y1) + x2 * (y1 - y3) + x1 * (y3 - y2)) * reversed_denom;
    *b = (x3 * x3 * (y1 - y2) + x2 * x2 * (y3 - y1) + x1 * x1 * (y2 - y3)) * reversed_denom;
    *c = (x2 * x3 * (x2 - x3) * y1 + x3 * x1 * (x3 - x1) * y2 + x1 * x2 * (x1 - x2) * y3) * reversed_denom;
}
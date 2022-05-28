#include "fft_module_o.h"
#include <Arduino.h>

byte sine_data[91] =
    {
        0,
        4, 9, 13, 18, 22, 27, 31, 35, 40, 44,
        49, 53, 57, 62, 66, 70, 75, 79, 83, 87,
        91, 96, 100, 104, 108, 112, 116, 120, 124, 127,
        131, 135, 139, 143, 146, 150, 153, 157, 160, 164,
        167, 171, 174, 177, 180, 183, 186, 189, 192, 195,
        198, 201, 204, 206, 209, 211, 214, 216, 219, 221,
        223, 225, 227, 229, 231, 233, 235, 236, 238, 240,
        241, 243, 244, 245, 246, 247, 248, 249, 250, 251,
        252, 253, 253, 254, 254, 254, 255, 255, 255, 255};

float Fft_Module_O::sine(int i)
{
    int j = i;
    float out;
    while (j < 0)
    {
        j = j + 360;
    }
    while (j > 360)
    {
        j = j - 360;
    }
    if (j > -1 && j < 91)
    {
        out = sine_data[j];
    }
    else if (j > 90 && j < 181)
    {
        out = sine_data[180 - j];
    }
    else if (j > 180 && j < 271)
    {
        out = -sine_data[j - 180];
    }
    else if (j > 270 && j < 361)
    {
        out = -sine_data[360 - j];
    }
    return (out / 255);
}

float Fft_Module_O::cosine(int i)
{
    int j = i;
    float out;
    while (j < 0)
    {
        j = j + 360;
    }
    while (j > 360)
    {
        j = j - 360;
    }
    if (j > -1 && j < 91)
    {
        out = sine_data[90 - j];
    }
    else if (j > 90 && j < 181)
    {
        out = -sine_data[j - 90];
    }
    else if (j > 180 && j < 271)
    {
        out = -sine_data[270 - j];
    }
    else if (j > 270 && j < 361)
    {
        out = sine_data[j - 270];
    }
    return (out / 255);
}

void Fft_Module_O::FFT(int in[], int N, float Frequency)
{
    unsigned int data[13] = {1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048};
    int a, c1, f, o, x;
    a = N;

    for (int i = 0; i < 12; i++) // calculating the levels
    {
        if (data[i] <= a)
        {
            o = i;
        }
    }

    int in_ps[data[o]] = {};
    float out_r[data[o]] = {};
    float out_im[data[o]] = {};

    x = 0;
    for (int b = 0; b < o; b++)
    {
        c1 = data[b];
        f = data[o] / (c1 + c1);
        for (int j = 0; j < c1; j++)
        {
            x = x + 1;
            in_ps[x] = in_ps[j] + f;
        }
    }

    for (int i = 0; i < data[o]; i++)
    {
        if (in_ps[i] < a)
        {
            out_r[i] = in[in_ps[i]];
        }
        if (in_ps[i] > a)
        {
            out_r[i] = in[in_ps[i] - a];
        }
    }

    int i10, i11, n1;
    float e, c, s, tr, ti;

    for (int i = 0; i < o; i++) // fft
    {
        i10 = data[i];
        i11 = data[o] / data[i + 1];
        e = 360 / data[i + 1];
        e = 0 - e;
        n1 = 0;

        for (int j = 0; j < i10; j++)
        {
            c = cosine(e * j);
            s = sine(e * j);
            n1 = j;

            for (int k = 0; k < i11; k++)
            {
                tr = c * out_r[i10 + n1] - s * out_im[i10 + n1];
                ti = s * out_r[i10 + n1] + c * out_im[i10 + n1];

                out_r[n1 + i10] = out_r[n1] - tr;
                out_r[n1] = out_r[n1] + tr;

                out_im[n1 + i10] = out_im[n1] - ti;
                out_im[n1] = out_im[n1] + ti;

                n1 = n1 + i10 + i10;
            }
        }
    }

    for (int i = 0; i < data[o - 1]; i++)
    {
        out_r[i] = sqrt(out_r[i] * out_r[i] + out_im[i] * out_im[i]);
        out_im[i] = i * Frequency / N;
    }

    x = 0;
    for (int i = 1; i < data[o - 1] - 1; i++)
    {
        if (out_r[i] > out_r[i - 1] && out_r[i] > out_r[i + 1])
        {
            in_ps[x] = i;
            x = x + 1;
        }
    }

    s = 0;
    c = 0;
    for (int i = 0; i < x; i++)
    {
        for (int j = c; j < x; j++)
        {
            if (out_r[in_ps[i]] < out_r[in_ps[j]])
            {
                s = in_ps[i];
                in_ps[i] = in_ps[j];
                in_ps[j] = s;
            }
        }
        c = c + 1;
    }

    for (int i = 0; i < 5; i++)
        f_peaks[i] = out_im[in_ps[i]];
}

void Fft_Module_O::perform_fft(float* current_fft, float* temperature_fft, float* vibration_fft, float* sound_fft)
{
    FFT(fft_input_data_temperature, MAX_SAMPLE_SIZE_A, 100);
    *temperature_fft = f_peaks[0];

    FFT(fft_input_data_current, MAX_SAMPLE_SIZE_A, 100);
    *current_fft = f_peaks[0];

    FFT(fft_input_data_vibration, MAX_SAMPLE_SIZE_A, 100);
    *vibration_fft = f_peaks[0];

    FFT(fft_input_data_sound, MAX_SAMPLE_SIZE_A, 100);
    *sound_fft = f_peaks[0];
}

void Fft_Module_O::recordSample(float current, float temperature, float vibration, float sound)
{
    fft_input_data_temperature[sample_counter] = temperature;
    fft_input_data_current[sample_counter] = current;
    fft_input_data_vibration[sample_counter] = vibration;
    fft_input_data_sound[sample_counter] = sound;
    sample_counter++;
}

int Fft_Module_O::getSampleCounterVal()
{
    return sample_counter;
}

void Fft_Module_O::resetSampleCounter()
{
    sample_counter = 0;
}

int Fft_Module_O::getMaxSampleSize()
{
    return MAX_SAMPLE_SIZE_A;
}
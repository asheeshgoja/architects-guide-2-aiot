#ifndef __FFT_MODULE__
#define __FFT_MODULE__

const int MAX_SAMPLE_SIZE_A = 64;

class Fft_Module_O
{
private:

    int fft_input_data_temperature[MAX_SAMPLE_SIZE_A] = {};
    int fft_input_data_current[MAX_SAMPLE_SIZE_A] = {};
    int fft_input_data_vibration[MAX_SAMPLE_SIZE_A] = {};
    int fft_input_data_sound[MAX_SAMPLE_SIZE_A] = {};
    int sample_counter = 0;

    float f_peaks[5];
    float sine(int i);
    float cosine(int i);
    void FFT(int in[], int N, float Frequency);

public:
    void perform_fft( float* current_fft, float* temperature_fft, float* vibration_fft, float* sound_fft);
    void recordSample(float current, float temperature, float vibration, float sound);
    int getSampleCounterVal();
    int getMaxSampleSize();
    void resetSampleCounter();
};

#endif

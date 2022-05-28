#ifndef ___TFLM_Module__
#define ___TFLM_Module__

#include <stdint.h>

// org at  ../kubecon_talk/esp32s
//  Globals, used for compatibility with Arduino-style sketches.
namespace tflite
{

    class ErrorReporter;
    class Model;
    class MicroInterpreter;
    template <unsigned int tOpCount>
    class MicroMutableOpResolver;

}

// forward decleration
struct TfLiteTensor;

class HTTPClient;
// class WiFiClass;

const int kTensorArenaSize = 20000;

class TFLM_Module
{

    int model_file_downloaded;

private:
    tflite::ErrorReporter *error_reporter;
    const tflite::Model *model;
    tflite::MicroInterpreter *interpreter;
    TfLiteTensor *input;
    TfLiteTensor *output;
    int inference_count;
    HTTPClient *pHTTPClient;
    // WiFiClass *pWiFi;
    char modelRegistry[255];

    uint8_t tensor_arena[kTensorArenaSize];
    int downloadTFLiteModel(uint8_t **tfLiteModel, const char *tfLiteFileName);
    float *getInputBuffer();

public:
    void init(HTTPClient *p_http, const char *reg);
    float predict(float current, float temperature, float vibration, float sound);
    void loadModel(const char *tfLiteFileName);
    static void setNewModelFileName(const char *tfLiteFileName);
};

#endif
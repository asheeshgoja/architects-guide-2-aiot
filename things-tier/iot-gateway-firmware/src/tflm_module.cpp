#include "tflm_module.h"
#include "tensorflow/lite/micro/all_ops_resolver.h"
#include "tensorflow/lite/micro/micro_error_reporter.h"
#include "tensorflow/lite/micro/micro_interpreter.h"
#include "tensorflow/lite/schema/schema_generated.h"
#include "tensorflow/lite/version.h"
#include <HTTPClient.h>


extern const char* reg_svr_pub_key;

char g_modelFileName[255] = {};
int g_newModelAvailable = 0;
int g_modelDownloadInProgress = 0;

void TFLM_Module::setNewModelFileName(const char *tfLiteFileName)
{
    g_newModelAvailable = 1;
    strcpy(g_modelFileName, tfLiteFileName);
    Serial.print("Model available for download - file name :  ");
    Serial.println(g_modelFileName);
}

int TFLM_Module::downloadTFLiteModel(uint8_t **tfLiteModel, const char *tfLiteFileName)
{

    char tfliteFileURL[255] = {};
    snprintf(tfliteFileURL, 255, "%s%s", modelRegistry, tfLiteFileName);

    Serial.printf("Downloading model file from registry server @%s\n", tfliteFileURL);
    // Your Domain name with URL path or IP address with path
    pHTTPClient->begin(tfliteFileURL, reg_svr_pub_key);

    // Send HTTP GET request
    int httpResponseCode = pHTTPClient->GET();

    if (httpResponseCode == 200)
    {
        Serial.print("HTTP GOOD Response code: ");
        Serial.println(httpResponseCode);

        // String payload = pHTTPClient->getString();
        // int len = payload.length()+1;
        int len = 11148;

        // Serial.print("TFLite model content: ");
        // Serial.println(payload);
        Serial.print("TFLite model bytes read: ");
        Serial.println(len);

        // payload = "ABCD";
        // len = payload.length() + 1;

        *tfLiteModel = (byte *)malloc(len);
        pHTTPClient->getStream().readBytes(*tfLiteModel, len);

        // payload.getBytes(*tfLiteModel, len);
        // printBinaryData(*tfLiteModel,len);

        return len;
        // Serial.println("Pay");
    }
    else
    {
        Serial.print("Error code: ");
        Serial.println(httpResponseCode);
        return -1;
    }
    // Free resources
    pHTTPClient->end();

    return -1;
}

void TFLM_Module::init(HTTPClient *p_http, const char *reg)
{
    pHTTPClient = p_http;
    model_file_downloaded = 0;
    strcpy(modelRegistry, reg);
    model_file_downloaded = 0;
}

void TFLM_Module::loadModel(const char *tfLiteFileName)
{
    uint8_t *tflileModelBuf = 0;

    Serial.print("loading model : ");
    Serial.println(tfLiteFileName);

    model_file_downloaded = downloadTFLiteModel(&tflileModelBuf, tfLiteFileName);

    if (model_file_downloaded)
    {
        // tflite::InitializeTarget();

        // Set up logging. Google style is to avoid globals or statics because of
        // lifetime uncertainty, but since this has a trivial destructor it's okay.
        // NOLINTNEXTLINE(runtime-global-variables)
        static tflite::MicroErrorReporter micro_error_reporter;
        error_reporter = &micro_error_reporter;

        // Map the model into a usable data structure. This doesn't involve any
        // copying or parsing, it's a very lightweight operation.
        model = tflite::GetModel(tflileModelBuf);
        if (model->version() != TFLITE_SCHEMA_VERSION)
        {
            TF_LITE_REPORT_ERROR(error_reporter,
                                 "Model provided is schema version %d not equal "
                                 "to supported version %d.",
                                 model->version(), TFLITE_SCHEMA_VERSION);
            return;
        }

        // This pulls in all the operation implementations we need.
        // NOLINTNEXTLINE(runtime-global-variables)
        static tflite::AllOpsResolver resolver;

        // Build an interpreter to run the model with.
        static tflite::MicroInterpreter static_interpreter(
            model, resolver, tensor_arena, kTensorArenaSize, error_reporter);
        interpreter = &static_interpreter;

        // Allocate memory from the tensor_arena for the model's tensors.
        TfLiteStatus allocate_status = interpreter->AllocateTensors();
        if (allocate_status != kTfLiteOk)
        {
            TF_LITE_REPORT_ERROR(error_reporter, "AllocateTensors() failed");
            return;
        }

        // Obtain pointers to the model's input and output tensors.
        input = interpreter->input(0);
        output = interpreter->output(0);

        // Keep track of how many inferences we have performed.
        inference_count = 0;
    }
}

float *TFLM_Module::getInputBuffer()
{
    return input->data.f;
    // return 0;
}

float TFLM_Module::predict(float current, float temperature, float vibration, float sound)
{
    // Serial.println("Model not available for download... ");

    if (g_newModelAvailable == 1)
    {
        g_modelDownloadInProgress = 1;
        loadModel(g_modelFileName);
        g_modelDownloadInProgress = 0;
        g_newModelAvailable = 0;
    }
    // else
    // {
    //     Serial.println("Model not available for download... ");
    //     return -1;
    // }

    if (model_file_downloaded)
    {

        getInputBuffer()[0] = current;
        getInputBuffer()[1] = temperature;
        getInputBuffer()[2] = vibration;
        getInputBuffer()[3] = sound;

        interpreter->Invoke();

        float ret_val = output->data.f[0];

        Serial.printf("TFLM inference val : %.2f, data : Current: %.2f, Temp: %.2f, Vibration: %.2f, Sound: %.2f\n", ret_val, current, temperature, vibration, sound);
    }
    return -1;
}

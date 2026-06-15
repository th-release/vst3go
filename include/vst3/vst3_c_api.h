#ifndef VST3GO_VST3_C_API_H
#define VST3GO_VST3_C_API_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef SMTG_STDMETHODCALLTYPE
#define SMTG_STDMETHODCALLTYPE
#endif

typedef int32_t Steinberg_tresult;
typedef int32_t Steinberg_int32;
typedef int64_t Steinberg_int64;
typedef uint32_t Steinberg_uint32;
typedef uint8_t Steinberg_TBool;
typedef uint8_t Steinberg_TUID[16];
typedef const char* Steinberg_FIDString;
typedef uint16_t Steinberg_char16;
typedef Steinberg_char16 Steinberg_Vst_TChar;
typedef Steinberg_Vst_TChar Steinberg_Vst_String128[128];
typedef double Steinberg_Vst_ParamValue;
typedef float Steinberg_Vst_Sample32;
typedef double Steinberg_Vst_Sample64;
typedef uint32_t Steinberg_Vst_ParamID;
typedef int32_t Steinberg_Vst_UnitID;
typedef uint64_t Steinberg_Vst_SpeakerArrangement;
typedef int32_t Steinberg_Vst_MediaType;
typedef int32_t Steinberg_Vst_BusDirection;
typedef int32_t Steinberg_Vst_BusType;
typedef int32_t Steinberg_Vst_IoMode;

enum {
    Steinberg_kResultOk = 0,
    Steinberg_kResultFalse = 1,
};

enum {
    Steinberg_PFactoryInfo_FactoryFlags_kUnicode = 1 << 0,
};

enum {
    Steinberg_PClassInfo_ClassCardinality_kManyInstances = 0x7fffffff
};

enum {
    Steinberg_Vst_MediaTypes_kAudio = 0,
    Steinberg_Vst_MediaTypes_kEvent = 1,
};

enum {
    Steinberg_Vst_BusDirections_kInput = 0,
    Steinberg_Vst_BusDirections_kOutput = 1,
};

enum {
    Steinberg_Vst_BusTypes_kMain = 0,
    Steinberg_Vst_BusTypes_kAux = 1,
};

enum {
    Steinberg_Vst_ParameterInfo_ParameterFlags_kCanAutomate = 1 << 0,
    Steinberg_Vst_ParameterInfo_ParameterFlags_kIsReadOnly = 1 << 1,
    Steinberg_Vst_ParameterInfo_ParameterFlags_kIsWrapAround = 1 << 2,
    Steinberg_Vst_ParameterInfo_ParameterFlags_kIsList = 1 << 3,
    Steinberg_Vst_ParameterInfo_ParameterFlags_kIsHidden = 1 << 4,
    Steinberg_Vst_ParameterInfo_ParameterFlags_kIsBypass = 1 << 16,
};

enum {
    Steinberg_Vst_Event_EventTypes_kNoteOffEvent = 0,
    Steinberg_Vst_Event_EventTypes_kNoteOnEvent = 1,
    Steinberg_Vst_Event_EventTypes_kPolyPressureEvent = 2,
    Steinberg_Vst_Event_EventTypes_kDataEvent = 3,
    Steinberg_Vst_Event_EventTypes_kLegacyMIDICCOutEvent = 4,
};

enum {
    Steinberg_Vst_ProcessContext_StatesAndFlags_kPlaying = 1 << 1,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kRecording = 1 << 2,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kCycleActive = 1 << 3,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kSystemTimeValid = 1 << 8,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kProjectTimeMusicValid = 1 << 9,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kTempoValid = 1 << 10,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kBarPositionValid = 1 << 11,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kCycleValid = 1 << 12,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kTimeSigValid = 1 << 13,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kSmpteValid = 1 << 14,
    Steinberg_Vst_ProcessContext_StatesAndFlags_kClockValid = 1 << 15,
};

typedef struct Steinberg_FUnknownVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
} Steinberg_FUnknownVtbl;

typedef struct Steinberg_FUnknown {
    Steinberg_FUnknownVtbl* lpVtbl;
} Steinberg_FUnknown;

typedef struct Steinberg_PFactoryInfo {
    char vendor[64];
    char url[256];
    char email[128];
    Steinberg_int32 flags;
} Steinberg_PFactoryInfo;

typedef struct Steinberg_PClassInfo {
    Steinberg_TUID cid;
    Steinberg_int32 cardinality;
    char category[64];
    char name[64];
} Steinberg_PClassInfo;

typedef struct Steinberg_IPluginFactoryVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getFactoryInfo)(void* thisInterface, Steinberg_PFactoryInfo* info);
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *countClasses)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getClassInfo)(void* thisInterface, Steinberg_int32 index, Steinberg_PClassInfo* info);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *createInstance)(void* thisInterface, Steinberg_FIDString cid, Steinberg_FIDString iid, void** obj);
} Steinberg_IPluginFactoryVtbl;

typedef struct Steinberg_IPluginFactory {
    Steinberg_IPluginFactoryVtbl* lpVtbl;
} Steinberg_IPluginFactory;

typedef struct Steinberg_IBStreamVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *read)(void* thisInterface, void* buffer, Steinberg_int32 numBytes, Steinberg_int32* numBytesRead);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *write)(void* thisInterface, void* buffer, Steinberg_int32 numBytes, Steinberg_int32* numBytesWritten);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *seek)(void* thisInterface, Steinberg_int64 pos, Steinberg_int32 mode, Steinberg_int64* result);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *tell)(void* thisInterface, Steinberg_int64* pos);
} Steinberg_IBStreamVtbl;

typedef struct Steinberg_IBStream {
    Steinberg_IBStreamVtbl* lpVtbl;
} Steinberg_IBStream;

typedef struct Steinberg_Vst_ProcessSetup {
    Steinberg_int32 processMode;
    Steinberg_int32 symbolicSampleSize;
    Steinberg_int32 maxSamplesPerBlock;
    double sampleRate;
} Steinberg_Vst_ProcessSetup;

typedef struct Steinberg_Vst_ProcessData Steinberg_Vst_ProcessData;

typedef struct Steinberg_Vst_ProcessContext {
    Steinberg_uint32 state;
    double sampleRate;
    Steinberg_int64 projectTimeSamples;
    Steinberg_int64 systemTime;
    Steinberg_int64 continousTimeSamples;
    double projectTimeMusic;
    double barPositionMusic;
    double cycleStartMusic;
    double cycleEndMusic;
    double tempo;
    Steinberg_int32 timeSigNumerator;
    Steinberg_int32 timeSigDenominator;
    Steinberg_int32 samplesToNextClock;
} Steinberg_Vst_ProcessContext;

typedef struct Steinberg_Vst_AudioBusBuffers {
    Steinberg_int32 numChannels;
    Steinberg_uint32 silenceFlags;
    union {
        Steinberg_Vst_Sample32** Steinberg_Vst_AudioBusBuffers_channelBuffers32;
        Steinberg_Vst_Sample64** Steinberg_Vst_AudioBusBuffers_channelBuffers64;
    };
} Steinberg_Vst_AudioBusBuffers;

typedef struct Steinberg_Vst_BusInfo {
    Steinberg_Vst_MediaType mediaType;
    Steinberg_Vst_BusDirection direction;
    Steinberg_int32 channelCount;
    Steinberg_Vst_TChar name[128];
    Steinberg_Vst_BusType busType;
    Steinberg_uint32 flags;
} Steinberg_Vst_BusInfo;

typedef struct Steinberg_Vst_RoutingInfo {
    Steinberg_Vst_MediaType mediaType;
    Steinberg_Vst_BusDirection direction;
    Steinberg_int32 busIndex;
    Steinberg_int32 channel;
    Steinberg_int32 channelOffset;
} Steinberg_Vst_RoutingInfo;

typedef struct Steinberg_Vst_ParameterInfo {
    Steinberg_Vst_ParamID id;
    Steinberg_Vst_TChar title[128];
    Steinberg_Vst_TChar shortTitle[128];
    Steinberg_Vst_TChar units[128];
    Steinberg_int32 stepCount;
    Steinberg_Vst_ParamValue defaultNormalizedValue;
    Steinberg_Vst_UnitID unitId;
    Steinberg_int32 flags;
} Steinberg_Vst_ParameterInfo;

typedef struct Steinberg_Vst_NoteOnEvent {
    Steinberg_int32 channel;
    Steinberg_int32 pitch;
    float tuning;
    Steinberg_Vst_ParamValue velocity;
    Steinberg_int32 length;
    Steinberg_int32 noteId;
} Steinberg_Vst_NoteOnEvent;

typedef struct Steinberg_Vst_NoteOffEvent {
    Steinberg_int32 channel;
    Steinberg_int32 pitch;
    float tuning;
    Steinberg_Vst_ParamValue velocity;
    Steinberg_int32 noteId;
} Steinberg_Vst_NoteOffEvent;

typedef struct Steinberg_Vst_Event {
    Steinberg_int32 type;
    Steinberg_int32 flags;
    Steinberg_int32 sampleOffset;
    union {
        Steinberg_Vst_NoteOnEvent Steinberg_Vst_Event_noteOn;
        Steinberg_Vst_NoteOffEvent Steinberg_Vst_Event_noteOff;
    };
} Steinberg_Vst_Event;

typedef struct Steinberg_Vst_IParamValueQueueVtbl {
    Steinberg_Vst_ParamID (SMTG_STDMETHODCALLTYPE *getParameterId)(void* thisInterface);
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *getPointCount)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getPoint)(void* thisInterface, Steinberg_int32 index, Steinberg_int32* sampleOffset, Steinberg_Vst_ParamValue* value);
} Steinberg_Vst_IParamValueQueueVtbl;

typedef struct Steinberg_Vst_IParamValueQueue {
    Steinberg_Vst_IParamValueQueueVtbl* lpVtbl;
} Steinberg_Vst_IParamValueQueue;

typedef struct Steinberg_Vst_IParameterChangesVtbl {
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *getParameterCount)(void* thisInterface);
    Steinberg_Vst_IParamValueQueue* (SMTG_STDMETHODCALLTYPE *getParameterData)(void* thisInterface, Steinberg_int32 index);
} Steinberg_Vst_IParameterChangesVtbl;

typedef struct Steinberg_Vst_IParameterChanges {
    Steinberg_Vst_IParameterChangesVtbl* lpVtbl;
} Steinberg_Vst_IParameterChanges;

typedef struct Steinberg_Vst_IEventListVtbl {
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *getEventCount)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getEvent)(void* thisInterface, Steinberg_int32 index, Steinberg_Vst_Event* event);
} Steinberg_Vst_IEventListVtbl;

typedef struct Steinberg_Vst_IEventList {
    Steinberg_Vst_IEventListVtbl* lpVtbl;
} Steinberg_Vst_IEventList;

typedef struct Steinberg_Vst_IComponentHandlerVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *beginEdit)(void* thisInterface, Steinberg_Vst_ParamID id);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *performEdit)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue value);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *endEdit)(void* thisInterface, Steinberg_Vst_ParamID id);
} Steinberg_Vst_IComponentHandlerVtbl;

typedef struct Steinberg_Vst_IComponentHandler {
    Steinberg_Vst_IComponentHandlerVtbl* lpVtbl;
} Steinberg_Vst_IComponentHandler;

typedef struct Steinberg_IPlugView {
    void* lpVtbl;
} Steinberg_IPlugView;

typedef struct Steinberg_Vst_IComponentVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *initialize)(void* thisInterface, Steinberg_FUnknown* context);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *terminate)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getControllerClassId)(void* thisInterface, Steinberg_TUID classId);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setIoMode)(void* thisInterface, Steinberg_Vst_IoMode mode);
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *getBusCount)(void* thisInterface, Steinberg_Vst_MediaType type, Steinberg_Vst_BusDirection dir);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getBusInfo)(void* thisInterface, Steinberg_Vst_MediaType type, Steinberg_Vst_BusDirection dir, Steinberg_int32 index, Steinberg_Vst_BusInfo* bus);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getRoutingInfo)(void* thisInterface, Steinberg_Vst_RoutingInfo* inInfo, Steinberg_Vst_RoutingInfo* outInfo);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *activateBus)(void* thisInterface, Steinberg_Vst_MediaType type, Steinberg_Vst_BusDirection dir, Steinberg_int32 index, Steinberg_TBool state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setActive)(void* thisInterface, Steinberg_TBool state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setState)(void* thisInterface, Steinberg_IBStream* state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getState)(void* thisInterface, Steinberg_IBStream* state);
} Steinberg_Vst_IComponentVtbl;

typedef struct Steinberg_Vst_IComponent {
    Steinberg_Vst_IComponentVtbl* lpVtbl;
} Steinberg_Vst_IComponent;

typedef struct Steinberg_Vst_IAudioProcessorVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setBusArrangements)(void* thisInterface, Steinberg_Vst_SpeakerArrangement* inputs, Steinberg_int32 numIns, Steinberg_Vst_SpeakerArrangement* outputs, Steinberg_int32 numOuts);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getBusArrangement)(void* thisInterface, Steinberg_Vst_BusDirection dir, Steinberg_int32 index, Steinberg_Vst_SpeakerArrangement* arr);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *canProcessSampleSize)(void* thisInterface, Steinberg_int32 symbolicSampleSize);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *getLatencySamples)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setupProcessing)(void* thisInterface, Steinberg_Vst_ProcessSetup* setup);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setProcessing)(void* thisInterface, Steinberg_TBool state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *process)(void* thisInterface, Steinberg_Vst_ProcessData* data);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *getTailSamples)(void* thisInterface);
} Steinberg_Vst_IAudioProcessorVtbl;

typedef struct Steinberg_Vst_IAudioProcessor {
    Steinberg_Vst_IAudioProcessorVtbl* lpVtbl;
} Steinberg_Vst_IAudioProcessor;

typedef struct Steinberg_Vst_IEditControllerVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *initialize)(void* thisInterface, Steinberg_FUnknown* context);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *terminate)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setComponentState)(void* thisInterface, Steinberg_IBStream* state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setState)(void* thisInterface, Steinberg_IBStream* state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getState)(void* thisInterface, Steinberg_IBStream* state);
    Steinberg_int32 (SMTG_STDMETHODCALLTYPE *getParameterCount)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getParameterInfo)(void* thisInterface, Steinberg_int32 paramIndex, Steinberg_Vst_ParameterInfo* info);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getParamStringByValue)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue valueNormalized, Steinberg_Vst_String128 string);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getParamValueByString)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_TChar* string, Steinberg_Vst_ParamValue* valueNormalized);
    Steinberg_Vst_ParamValue (SMTG_STDMETHODCALLTYPE *normalizedParamToPlain)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue valueNormalized);
    Steinberg_Vst_ParamValue (SMTG_STDMETHODCALLTYPE *plainParamToNormalized)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue plainValue);
    Steinberg_Vst_ParamValue (SMTG_STDMETHODCALLTYPE *getParamNormalized)(void* thisInterface, Steinberg_Vst_ParamID id);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setParamNormalized)(void* thisInterface, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue value);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setComponentHandler)(void* thisInterface, Steinberg_Vst_IComponentHandler* handler);
    Steinberg_IPlugView* (SMTG_STDMETHODCALLTYPE *createView)(void* thisInterface, Steinberg_FIDString name);
} Steinberg_Vst_IEditControllerVtbl;

typedef struct Steinberg_Vst_IEditController {
    Steinberg_Vst_IEditControllerVtbl* lpVtbl;
} Steinberg_Vst_IEditController;

typedef struct Steinberg_Vst_ProcessData {
    Steinberg_int32 processMode;
    Steinberg_int32 symbolicSampleSize;
    Steinberg_int32 numSamples;
    Steinberg_int32 numInputs;
    Steinberg_int32 numOutputs;
    Steinberg_Vst_AudioBusBuffers* inputs;
    Steinberg_Vst_AudioBusBuffers* outputs;
    Steinberg_Vst_IEventList* inputEvents;
    Steinberg_Vst_IEventList* outputEvents;
    Steinberg_Vst_IParameterChanges* inputParameterChanges;
    Steinberg_Vst_IParameterChanges* outputParameterChanges;
    Steinberg_Vst_ProcessContext* processContext;
} Steinberg_Vst_ProcessData;

static const Steinberg_TUID Steinberg_FUnknown_iid = {
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46
};

static const Steinberg_TUID Steinberg_IPluginBase_iid = {
    0xEF, 0xBC, 0x16, 0x64, 0x15, 0xD4, 0x4A, 0x2D,
    0x83, 0xA9, 0xB2, 0x10, 0x43, 0x9B, 0x00, 0x00
};

static const Steinberg_TUID Steinberg_IPluginFactory_iid = {
    0x7A, 0x4D, 0x81, 0x1C, 0x52, 0x11, 0x4A, 0x1F,
    0xAE, 0xD9, 0xD2, 0xEE, 0x0B, 0x43, 0xBF, 0x9F
};

static const Steinberg_TUID Steinberg_Vst_IComponent_iid = {
    0xE8, 0x31, 0xFF, 0x31, 0xF2, 0xD5, 0x4B, 0x01,
    0x83, 0x6F, 0x5D, 0x38, 0x54, 0x34, 0xAE, 0xC6
};

static const Steinberg_TUID Steinberg_Vst_IAudioProcessor_iid = {
    0x42, 0x04, 0x3F, 0x99, 0xB2, 0xA8, 0x4F, 0x3F,
    0xA2, 0x85, 0x7A, 0xA0, 0x39, 0x82, 0x15, 0xC1
};

static const Steinberg_TUID Steinberg_Vst_IEditController_iid = {
    0xDD, 0xB1, 0x18, 0x8F, 0x2B, 0x0D, 0x43, 0x11,
    0x9E, 0xD0, 0xAE, 0xB4, 0x38, 0x95, 0x40, 0x52
};

#ifdef __cplusplus
}
#endif

#endif

package desktopCapture

// #cgo LDFLAGS: -lavdevice -lavcodec -lavformat -lavutil -lswscale -lx265 -lSDL2
/*
#define __STDC_CONSTANT_MACROS
#define SDL_MAIN_HANDLED

#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif

#include <libavutil/avutil.h>
#include <libavutil/imgutils.h>
#include <libswscale/swscale.h>
#include <libavcodec/avcodec.h>
#include <libavdevice/avdevice.h>
#include <libavfilter/avfilter.h>
#include <libavformat/avformat.h>
#include <libswresample/swresample.h>
#include <x265.h>
#include <SDL2/SDL.h>

#ifdef __cplusplus
}
#endif

//#define USE_X264
#define USE_X265

#define ERR_AVFORMAT_ALLOC_CONTEXT -1
#define ERR_AV_FIND_INPUT_FORMAT -2
#define ERR_AVFORMAT_OPEN_INPUT -3
#define ERR_AVFORMAT_FIND_STREAM_INFO -4
#define ERR_AVMEDIA_TYPE_VIDEO_NOT_FOUND -5
#define ERR_AVCODEC_FIND_ENCODER -6
#define ERR_AVCODEC_ALLOC_CONTEXT3 -7
#define ERR_AVCODEC_PARAMETERS_TO_CONTEXT -8
#define ERR_AVCODEC_OPEN2 -9
#define ERR_X265_ENCODER_OPEN -10
#define ERR_AV_READ_FRAME -11
#define ERR_CSP_NOT_SUPPORTED -12
#define ERR_SDL_CREATE_WINDOW -13
#define ERR_SDL_CREATE_RENDERER -14
#define ERR_SDL_CREATE_TEXTURE - 15

#if defined(__clang__)
#define INPUT_FORMAT "" //TODO
#elif defined(__GNUC__) || defined(__GNUG__)
#ifdef __MINGW32__
#define INPUT_FORMAT "gdigrab"
#define INPUT_URL "desktop"
#else
#define INPUT_URL ":0.0"
#define INPUT_FORMAT "x11grab"
#endif
#elif defined(_MSC_VER)
#define INPUT_FORMAT "gdigrab"
#define INPUT_URL "desktop"
#endif

// show av module info
#define SHOW_MODULE_INFO(MODULE) \
    printf("(module: %s)(version: %d)(configuration: %s)(license: %s)\n", #MODULE, \
    MODULE##_version(),          \
    MODULE##_configuration(),    \
    MODULE##_license())

typedef void(*OnX265Nal)(uint32_t type,
                         uint32_t sizeBytes,
                         uint8_t *payload);

// DesktopCapture
struct DesktopCapture {
  int (*create)(struct DesktopCapture *dc);

  void (*destroy)(struct DesktopCapture *dc);

  void (*shoErrorMessage)(int ec);

  int (*x265Setup)(struct DesktopCapture *dc);

  int (*x265Encode)(struct DesktopCapture *ds,
                    struct AVFrame *frameYUV,
                    int internalCsp);

  int (*x265Flush)(struct DesktopCapture *dc);

  void (*onX265Nal)(uint32_t type,
                    uint32_t sizeBytes,
                    uint8_t *payload);

  int (*run)(struct DesktopCapture *dc);

  void (*onYUVFrame)(struct DesktopCapture *dc, struct AVFrame *frame);

  void (*stop)(struct DesktopCapture *dc);

  struct DesktopCapture *(*setFrameRate)(struct DesktopCapture *dc, int frameRate);

  struct DesktopCapture *(*setSrcOffsetX)(struct DesktopCapture *dc, int srcOffsetX);

  struct DesktopCapture *(*setSrcOffsetY)(struct DesktopCapture *dc, int srcOffsetY);

  struct DesktopCapture *(*setSrcWidth)(struct DesktopCapture *dc, int srcWidth);

  struct DesktopCapture *(*setSrcHeight)(struct DesktopCapture *dc, int srcHeight);

  struct DesktopCapture *(*setDstWidth)(struct DesktopCapture *dc, int dstWidth);

  struct DesktopCapture *(*setDstHeight)(struct DesktopCapture *dc, int dstHeight);

  struct DesktopCapture *(*setDstPixelFmt)(struct DesktopCapture *dc, enum AVPixelFormat dstPixelFmt);

  struct DesktopCapture *(*setOnX265Nal)(struct DesktopCapture *dc, OnX265Nal onX265Nal);

  int _stop;
  int _frameRate;
  int _srcOffsetX;
  int _srcOffsetY;
  int _srcWidth;
  int _srcHeight;
  int _dstWidth;
  int _dstHeight;
  Uint32 _updatePreviewEvent;
  enum AVPixelFormat _dstPixelFmt;
  OnX265Nal _onX265Nal;
  unsigned int _videoIndex;
  struct AVFormatContext *_fmtCtx;
  struct x265_encoder *_x265Encoder;
  struct x265_picture *_x265Picture;
  struct x265_nal *_x265Nal;
  const char *_windowTitle;
  struct SDL_Window *_previewWindow;
  struct SDL_Renderer *_previewRenderer;
  struct SDL_Texture *_previewTexture;
};

// create
int _create(struct DesktopCapture *dc) {
  int ec = 0;
  SHOW_MODULE_INFO(avdevice);
  SHOW_MODULE_INFO(avcodec);
  SHOW_MODULE_INFO(avformat);
  SHOW_MODULE_INFO(avutil);
  dc->_stop = 0;
  dc->_videoIndex = -1;
  dc->_dstWidth = 0;
  dc->_dstHeight = 0;
  dc->_frameRate = 0;
  dc->_updatePreviewEvent = 0;
  dc->_dstPixelFmt = AV_PIX_FMT_YUV420P;
  dc->_x265Nal = NULL;
  dc->_x265Encoder = NULL;
  dc->_x265Picture = NULL;
  dc->onX265Nal = NULL;
  dc->_windowTitle = "desktop preview";
  dc->_previewWindow = NULL;
  dc->_previewRenderer = NULL;
  dc->_previewTexture = NULL;
  dc->_fmtCtx = avformat_alloc_context();
  if (NULL == dc->_fmtCtx) return ERR_AVFORMAT_ALLOC_CONTEXT;
  if (ec = SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO | SDL_INIT_EVENTS)) {
    fprintf(stderr, "SDL_Init error: %s\n", SDL_GetError());
    return ec;
  }
  avdevice_register_all();
  return 0;
}

// destrory
void _destroy(struct DesktopCapture *dc) {
  if (dc->_fmtCtx) avformat_free_context(dc->_fmtCtx);
  if (dc->_x265Picture) x265_picture_free(dc->_x265Picture);
  if (dc->_x265Encoder) x265_encoder_close(dc->_x265Encoder);
  x265_cleanup();
  if (dc->_previewTexture) SDL_DestroyTexture(dc->_previewTexture);
  if (dc->_previewRenderer) SDL_DestroyRenderer(dc->_previewRenderer);
  if (dc->_previewWindow) SDL_DestroyWindow(dc->_previewWindow);
  SDL_Quit();
}

// show error message
void _showErrorMessage(int ec) {
  char buf[AV_ERROR_MAX_STRING_SIZE];
  av_make_error_string(buf, AV_ERROR_MAX_STRING_SIZE, ec);
  printf("%s\n", buf);
}

// set frame rate
struct DesktopCapture *_setFrameRate(struct DesktopCapture *dc, int frameRate) {
  dc->_frameRate = frameRate;
  return dc;
}

// set src offset x
struct DesktopCapture *_setSrcOffsetX(struct DesktopCapture *dc, int srcOffsetX) {
  dc->_srcOffsetX = srcOffsetX;
  return dc;
}

// set src offset y
struct DesktopCapture *_setSrcOffsetY(struct DesktopCapture *dc, int srcOffsetY) {
  dc->_srcOffsetY = srcOffsetY;
  return dc;
}

// set src width
struct DesktopCapture *_setSrcWidth(struct DesktopCapture *dc, int srcWidth) {
  dc->_srcWidth = srcWidth;
  return dc;
}

// set src height
struct DesktopCapture *_setSrcHeight(struct DesktopCapture *dc, int srcHeight) {
  dc->_srcHeight = srcHeight;
  return dc;
}

// set dst width
struct DesktopCapture *_setDstWidth(struct DesktopCapture *dc, int dstWidth) {
  dc->_dstWidth = dstWidth;
  return dc;
}

// set dst height
struct DesktopCapture *_setDstHeight(struct DesktopCapture *dc, int dstHeight) {
  dc->_dstHeight = dstHeight;
  return dc;
}

// set dst pixel format
struct DesktopCapture *_setDstPixelFmt(struct DesktopCapture *dc, enum AVPixelFormat dstPixelFmt) {
  dc->_dstPixelFmt = dstPixelFmt;
  return dc;
}

// set x265 nal
struct DesktopCapture *_setOnX265Nal(struct DesktopCapture *dc, OnX265Nal onX265Nal) {
  dc->_onX265Nal = onX265Nal;
  return dc;
}

// x265setup
int _x265setup(struct DesktopCapture *dc) {
  struct x265_param *param = x265_param_alloc();
  x265_param_default(param);
  param->bRepeatHeaders = 1;
  param->internalCsp = X265_CSP_I420;
  param->sourceWidth = dc->_dstWidth;
  param->sourceHeight = dc->_dstHeight;
  param->fpsNum = dc->_frameRate;
  param->fpsDenom = 1;
  param->logLevel = X265_LOG_NONE;
  param->frameNumThreads = 1;  // 限制CPU使用率
  // auto ret = x265_param_apply_profile(x265Param, x265_profile_names[1]);
  // if (0 > ret) LOG(ERROR) << "x265_param_apply_profile error:" << ret;
  dc->_x265Encoder = x265_encoder_open(param);
  if (!dc->_x265Encoder) return ERR_X265_ENCODER_OPEN;
  dc->_x265Picture = x265_picture_alloc();
  x265_picture_init(param, dc->_x265Picture);
  x265_param_free(param);
  return 0;
}

// _x265Encode
int _x265Encode(struct DesktopCapture *dc,
                struct AVFrame *frameYUV,
                int internalCsp) {
  int ec = 0;
  unsigned int nalNumber = 0;
  unsigned int i = 0;
  struct x265_picture *picture = dc->_x265Picture;
  switch (internalCsp) {
    case X265_CSP_I420: {
// C11
//      picture->planes[0] = frameYUV->data[0];
//      picture->planes[1] = frameYUV->data[1];
//      picture->planes[2] = frameYUV->data[2];
//      picture->stride[0] = width;
//      picture->stride[1] = width / 2;
//      picture->stride[2] = width / 2;
      // C99
      *(picture->planes) = frameYUV->data[0];
      *(picture->planes + 1) = frameYUV->data[1];
      *(picture->planes + 2) = frameYUV->data[2];
      *(picture->stride) = dc->_dstWidth;
      *(picture->stride + 1) = dc->_dstWidth / 2;
      *(picture->stride + 2) = dc->_dstWidth / 2;
    }
      break;
    case X265_CSP_I444: {
    }
      break;
    default: {
      return ERR_CSP_NOT_SUPPORTED;
    }
  }

  // 获取SPS&&PPS用于初始化解码器，在生命周期内可保持唯一，x265Param->bRepeatHeaders可设置为0，让不在产生IDR帧
  // int x265_encoder_headers(x265_encoder*, x265_nal * *pp_nal,
  //                         uint32_t * pi_nal);

  ec = x265_encoder_encode(dc->_x265Encoder, &dc->_x265Nal, &nalNumber,
                           dc->_x265Picture, NULL);
  if (0 > ec) {
    return ec;
  }
  for (i = 0; i < nalNumber; i++)
    if (dc->_onX265Nal)
      dc->_onX265Nal(dc->_x265Nal[i].type,
                    dc->_x265Nal[i].sizeBytes,
                    dc->_x265Nal[i].payload);
  return 0;
}

// _x265Flush
int _x265Flush(struct DesktopCapture *dc) {
  int ec = 0;
  unsigned int nalNumber = 0;
  unsigned int i = 0;
  while (1) {
    if (0 == (ec = x265_encoder_encode(dc->_x265Encoder,
                                       &dc->_x265Nal,
                                       &nalNumber,
                                       NULL,
                                       NULL)))
      break;
    if (0 > ec)
      return ec;
    for (i = 0; i < nalNumber; i++)
      if (dc->_onX265Nal)
        dc->_onX265Nal(dc->_x265Nal[i].type,
                      dc->_x265Nal[i].sizeBytes,
                      dc->_x265Nal[i].payload);
  }
  return 0;
}

// capture
int _capture(void *arg) {
  //
  struct DesktopCapture *dc = (struct DesktopCapture *) arg;
  int ec = 0;
  AVInputFormat *inputFmt = NULL;
  char value[32];
  const char *inputURL = NULL;
  AVDictionary *dict = NULL;
  int i = 0;
  AVCodecParameters *codecpar = NULL;
  AVCodec *decodec = NULL;
  struct AVCodecContext *ctx3 = NULL;
  struct AVFrame *frameYUV = NULL;
  struct SwsContext *swsCtx = NULL;
  int height = 0;
  struct AVPacket *packet = NULL;
  struct AVFrame *frameDesktop = NULL;
  do {
    // find input format
    printf("INPUT_FORMAT:%s\n", INPUT_FORMAT);
    if (NULL == (inputFmt = av_find_input_format(INPUT_FORMAT))) {
      ec = ERR_AV_FIND_INPUT_FORMAT;
      break;
    }
    // set frame rate
    av_dict_set_int(&dict, "framerate", dc->_frameRate, AV_DICT_MATCH_CASE);
    // set src offset x
    av_dict_set_int(&dict, "offset_x", dc->_srcOffsetX, AV_DICT_MATCH_CASE);
    // set src offset y
    av_dict_set_int(&dict, "offset_y", dc->_srcOffsetY, AV_DICT_MATCH_CASE);
    // set src width
    value[sprintf(value, "%dx%d", dc->_srcWidth, dc->_srcHeight)] = 0;
    av_dict_set(&dict, "video_size", value, 0);
    inputURL = INPUT_URL;
    for (i = 0; i < 2; i++) {
      ec = 0;
      if (0 > (ec = avformat_open_input(&dc->_fmtCtx, inputURL, inputFmt, &dict))) {
        dc->shoErrorMessage(AVERROR(ec));
        ec = ERR_AVFORMAT_OPEN_INPUT;
#if defined(__GNUC__) || defined(__GNUG__)
        inputURL = getenv("DISPLAY");
        printf("display: %s\n", inputURL);
#else
        break;
#endif
      } else break;
    }
    if (ec) return ec;
    if (0 > (ec = avformat_find_stream_info(dc->_fmtCtx, NULL))) {
      dc->shoErrorMessage(AVERROR(ec));
      ec = ERR_AVFORMAT_FIND_STREAM_INFO;
      break;
    }
    for (i = 0; i < dc->_fmtCtx->nb_streams; i++) {
      if (AVMEDIA_TYPE_VIDEO == dc->_fmtCtx->streams[i]->codecpar->codec_type) {
        dc->_videoIndex = i;
        break;
      }
    }
    if (0 > dc->_videoIndex) {
      ec = ERR_AVMEDIA_TYPE_VIDEO_NOT_FOUND;
      break;
    }
    codecpar = dc->_fmtCtx->streams[dc->_videoIndex]->codecpar;
    decodec = avcodec_find_decoder(codecpar->codec_id);
    if (!decodec) {
      ec = ERR_AVCODEC_FIND_ENCODER;
      break;
    }
    ctx3 = avcodec_alloc_context3(decodec);
    if (!ctx3) {
      ec = ERR_AVCODEC_ALLOC_CONTEXT3;
      break;
    }
    if (0 > (ec = avcodec_parameters_to_context(ctx3, codecpar))) {
      dc->shoErrorMessage(AVERROR(ec));
      ec = ERR_AVCODEC_PARAMETERS_TO_CONTEXT;
      break;
    }
    if (0 > (ec = avcodec_open2(ctx3, decodec, NULL))) {
      dc->shoErrorMessage(AVERROR(ec));
      ec = ERR_AVCODEC_OPEN2;
      break;
    }
    printf("(%d,%d,%d,%d,%d,%d,%d)\n",
           dc->_srcOffsetX,
           dc->_srcOffsetY,
           dc->_srcWidth,
           dc->_srcHeight,
           dc->_dstWidth,
           dc->_dstHeight,
           dc->_dstPixelFmt);
    frameYUV = av_frame_alloc();
    av_image_fill_arrays(
        frameYUV->data, frameYUV->linesize,
        (uint8_t *) av_malloc(av_image_get_buffer_size(dc->_dstPixelFmt,
                                                       dc->_dstWidth,
                                                       dc->_dstHeight,
                                                       1) * sizeof(uint8_t)),
        dc->_dstPixelFmt,
        dc->_dstWidth,
        dc->_dstHeight,
        1);
    swsCtx = sws_getContext(dc->_srcWidth,
                            dc->_srcHeight,
                            ctx3->pix_fmt,
                            dc->_dstWidth,
                            dc->_dstHeight,
                            dc->_dstPixelFmt,
                            SWS_BICUBIC,
                            NULL,
                            NULL,
                            NULL);
    packet = (struct AVPacket *) av_malloc(sizeof(struct AVPacket));
    frameDesktop = av_frame_alloc();
#if defined(USE_X264)
    //TODO
#elif defined(USE_X265)
    if (0 != (ec = dc->x265Setup(dc))) break;
#endif

    while (!dc->_stop) {
      if (0 > (ec = av_read_frame(dc->_fmtCtx, packet))) {
        dc->shoErrorMessage(AVERROR(ec));
        ec = ERR_AV_READ_FRAME;
        break;
      }
      if (0 == avcodec_send_packet(ctx3, packet)) {
        while (0 == avcodec_receive_frame(ctx3, frameDesktop)) {
          height = sws_scale(swsCtx,
                             frameDesktop->data,
                             frameDesktop->linesize,
                             0,
                             frameDesktop->height,
                             frameYUV->data,
                             frameYUV->linesize);
#if defined(USE_X264)
          //TODO
#elif defined(USE_X265)
          dc->onYUVFrame(dc, frameYUV);
          dc->x265Encode(dc, frameYUV, X265_CSP_I420);
#endif
        }
        if (0 != ec && AVERROR(EAGAIN) != ec)
          dc->shoErrorMessage(AVERROR(ec));
      } else {
        if (AVERROR_EOF == ec)
          dc->_stop = 1;
        else if (AVERROR(EAGAIN) == ec)
          continue;
        else if (0 != ec)
          dc->shoErrorMessage(AVERROR(ec));
      }
      av_packet_unref(packet);
    }
#if defined(USE_X264)
    //TODO
#elif defined(USE_X265)
    dc->x265Flush(dc);
#endif
    if (frameDesktop) av_frame_free(&frameDesktop);
    if (frameYUV) av_frame_free(&frameYUV);
    if (packet) av_free(packet);
  } while (0);
  return ec;
}

// run
int _run(struct DesktopCapture *dc) {
  int ec = 0;
  int runningFlag = 1;
  struct SDL_Thread *captureThread = NULL;
  union SDL_Event e;
  do {
    // show preview window
    if (!(dc->_previewWindow = SDL_CreateWindow(dc->_windowTitle,
                                                SDL_WINDOWPOS_CENTERED,
                                                SDL_WINDOWPOS_CENTERED,
                                                dc->_dstWidth,
                                                dc->_dstHeight,
                                                SDL_WINDOW_OPENGL))) {
      fprintf(stderr, "SDL_CreateWindow error: %s\n", SDL_GetError());
      ec = ERR_SDL_CREATE_WINDOW;
      break;
    }
    if (!(dc->_previewRenderer = SDL_CreateRenderer(dc->_previewWindow,
                                                    -1,
                                                    SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC))) {
      fprintf(stderr, "SDL_CreateRenderer error: %s\n", SDL_GetError());
      ec = ERR_SDL_CREATE_RENDERER;
      break;
    }
    if (!(dc->_previewTexture = SDL_CreateTexture(dc->_previewRenderer,
                                                  SDL_PIXELFORMAT_YV12,
                                                  SDL_TEXTUREACCESS_TARGET,
                                                  dc->_dstWidth,
                                                  dc->_dstHeight))) {
      fprintf(stderr, "SDL_CreateTexture error: %s\n", SDL_GetError());
      ec = ERR_SDL_CREATE_TEXTURE;
      break;
    }
    captureThread = SDL_CreateThread(_capture, "capture", dc);
    while (runningFlag) {
      if (SDL_PollEvent(&e)) {
        if (SDL_QUIT == e.type) {
          fprintf(stderr, "SDL_QUIT\n");
          ((struct DesktopCapture *) dc)->_stop = 1;
          runningFlag = 0;
        } else if (dc->_updatePreviewEvent == e.type) {
          struct AVFrame *frame = (struct AVFrame *) e.user.data1;
          if (0 > SDL_UpdateYUVTexture(dc->_previewTexture,
                                       NULL,
                                       frame->data[0],
                                       frame->linesize[0],
                                       frame->data[1],
                                       frame->linesize[1],
                                       frame->data[2],
                                       frame->linesize[2])) {
            fprintf(stderr, "SDL_UpdateYUVTexture error: %s\n", SDL_GetError());
          } else {
            SDL_RenderClear(dc->_previewRenderer);
            SDL_RenderCopy(dc->_previewRenderer, dc->_previewTexture, NULL, NULL);
            SDL_RenderPresent(dc->_previewRenderer);
          }
          free(frame->data[0]);
          free(frame->data[1]);
          free(frame->data[2]);
          free(frame);
        }
      }
    }
    dc->_stop = 1;
    SDL_WaitThread(captureThread, &ec);
  } while (0);
  return ec;
}

// onYUVFrame
void _onYUVFrame(struct DesktopCapture *dc, struct AVFrame *frame) {
  struct AVFrame *copied = (struct AVFrame *) malloc(sizeof(struct AVFrame));
  int lineSize = dc->_dstWidth * dc->_dstHeight;
  SDL_Event e;
  copied->data[0] = malloc(lineSize * sizeof(char));
  copied->data[1] = malloc(lineSize / 4 * sizeof(char));
  copied->data[2] = malloc(lineSize / 4 * sizeof(char));
  copied->linesize[0] = frame->linesize[0];
  copied->linesize[1] = frame->linesize[1];
  copied->linesize[2] = frame->linesize[2];
  memcpy(copied->data[0], frame->data[0], lineSize);
  memcpy(copied->data[1], frame->data[1], lineSize / 4);
  memcpy(copied->data[2], frame->data[2], lineSize / 4);
  dc->_updatePreviewEvent = SDL_RegisterEvents(1);
  SDL_zero(e);
  e.type = dc->_updatePreviewEvent;
  e.user.data1 = copied;
  SDL_PushEvent(&e);
}

// stop
void _stop(struct DesktopCapture *dc) {
  if (1 == dc->_stop) return;
  union SDL_Event e;
  e.type = SDL_QUIT;
  SDL_PushEvent(&e);
  dc->_stop = 1;
}

// makeDesktopCapture
struct DesktopCapture makeDesktopCapture() {
  struct DesktopCapture dc;
  dc.create = _create;
  dc.destroy = _destroy;
  dc.shoErrorMessage = _showErrorMessage;
  dc.x265Setup = _x265setup;
  dc.x265Encode = _x265Encode;
  dc.x265Flush = _x265Flush;
  dc.run = _run;
  dc.stop = _stop;
  dc.onYUVFrame = _onYUVFrame;
  dc.setFrameRate = _setFrameRate;
  dc.setSrcOffsetX = _setSrcOffsetX;
  dc.setSrcOffsetY = _setSrcOffsetY;
  dc.setSrcWidth = _setSrcWidth;
  dc.setSrcHeight = _setSrcHeight;
  dc.setDstWidth = _setDstWidth;
  dc.setDstHeight = _setDstHeight;
  dc.setDstPixelFmt = _setDstPixelFmt;
  dc.setOnX265Nal = _setOnX265Nal;
  return dc;
}

// desktop capture
struct DesktopCapture dc;

// startCapture
int startCapture(int frameRate,
                 int srcOffsetX,
                 int srcOffsetY,
                 int srcWidth,
                 int srcHeight,
                 int dstWidth,
                 int dstHeight,
                 int dstPixelFmt) {
  extern void onX265Nal(uint32_t type, uint32_t sizeBytes, uint8_t *payload);
  int ec;
  dc = makeDesktopCapture();
  if (0 != (ec = dc.create(&dc)))
    return ec;
  ec = dc.setFrameRate(&dc, frameRate)
      ->setSrcOffsetX(&dc, srcOffsetX)
      ->setSrcOffsetY(&dc, srcOffsetY)
      ->setSrcWidth(&dc, srcWidth)
      ->setSrcHeight(&dc, srcHeight)
      ->setDstWidth(&dc, dstWidth)
      ->setDstHeight(&dc, dstHeight)
      ->setDstPixelFmt(&dc, (enum AVPixelFormat)dstPixelFmt)
      ->setOnX265Nal(&dc, onX265Nal)
      ->run(&dc);
  dc.destroy(&dc);
  return ec;
}

// stopCapture
void stopCapture() {
  dc.stop(&dc);
}
*/
import "C"

func StartCapture(frameRate,
	srcOffsetX,
	srcOffsetY,
	srcWidth,
	srcHeight,
	dstWidth,
	dstHeight,
	dstPixelFmt int) int {
	return int(C.startCapture(C.int(frameRate),
		C.int(srcOffsetX),
		C.int(srcOffsetY),
		C.int(srcWidth),
		C.int(srcHeight),
		C.int(dstWidth),
		C.int(dstHeight),
		C.int(dstPixelFmt)))
}

func StopCapture() {
	C.stopCapture()
}

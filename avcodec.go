package sfutils

import "C"
import (
	"unsafe"
)

// #cgo CFLAGS: -std=c99
// #cgo LDFLAGS: -lavcodec -lavformat -lavutil
/*
#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>

typedef struct mediaInfo {
    int ok;
    char format[32];
    int width;
    int height;
    double duration;
} mediaInfo;

typedef struct mediaOpaque
{
    uint8_t *data;
    size_t length;
    size_t read;
} mediaOpaque;

static int mediaRead(void *pointer, uint8_t *buffer, int buf_size)
{
    mediaOpaque *opaque = (mediaOpaque *)pointer;
    size_t needed = buf_size;

    if(opaque->length - opaque->read < needed)
    {
        needed = opaque->length - opaque->read;
    }

    memcpy(buffer, opaque->data + opaque->read, needed);
    opaque->read += needed;

    return needed;
}

void getMediaDuration(mediaInfo *result, AVFormatContext *format, int videoStream, int audioStream) {
    AVStream *stream = NULL;
    if (videoStream != AVERROR_STREAM_NOT_FOUND) {
        stream = format->streams[videoStream];
    } else if (audioStream != AVERROR_STREAM_NOT_FOUND) {
        stream = format->streams[audioStream];
    }
    if (stream == NULL) {
        return;
    }

    result->duration = (double)stream->duration * av_q2d(stream->time_base);

    // Still, duration can be fucked, because some files don't provide it and you must calculate duration based on bitrate.
}

mediaInfo getMediaInfo(uint8_t *data, size_t length) {
    AVIOContext *io;
    AVFormatContext *format;
    mediaOpaque opaque;
    uint8_t *buffer = (uint8_t *)av_malloc(4096);

    mediaInfo result;
    result.ok = 0;
    result.format[0] = '\0';
    result.width = 0;
    result.height = 0;
    result.duration = 0.0;

    static int initialized = 0;
    if(initialized == 0) {
#ifndef FF_API_NEXT
        av_register_all();
#endif
        av_log_set_level(AV_LOG_QUIET);
    }

    opaque.data = data;
    opaque.length = length;
    opaque.read = 0;

    io = avio_alloc_context(buffer, 4096, 0, &opaque, &mediaRead, NULL, NULL);
    format = avformat_alloc_context();
    format->pb = io;

    // load file header
    if(avformat_open_input(&format, "buffer", NULL, NULL) >= 0 && format != NULL) {
        if(avformat_find_stream_info(format, NULL) >= 0) {
            // find video stream
            if(strstr(format->iformat->name, "hls") == NULL &&
               strstr(format->iformat->name, "http") == NULL) {
                strncpy(result.format, format->iformat->name, 32);
                result.format[31] = '\0';
                result.ok = 1;
            }

            if(result.ok) {
                int videoStream = av_find_best_stream(format, AVMEDIA_TYPE_VIDEO, -1, -1, NULL, 0);
                int audioStream = av_find_best_stream(format, AVMEDIA_TYPE_AUDIO, -1, -1, NULL, 0);

                if(videoStream != AVERROR_STREAM_NOT_FOUND) {
                    result.width = format->streams[videoStream]->codecpar->width;
                    result.height = format->streams[videoStream]->codecpar->height;
                }

                result.duration = (double)(format->duration) / AV_TIME_BASE;
                //getMediaDuration(&result, format, videoStream, audioStream);
            }
        }

        avformat_close_input(&format);
    }

    av_free(io->buffer);
#ifndef FF_API_NEXT
    av_freep(&io);
#else
    avio_context_free(&io);
#endif

    return result;
}*/
import "C"

type MediaInfo struct {
	Format   string
	Width    int
	Height   int
	Duration float64
}

func GetMediaInfo(data []byte) (*MediaInfo, bool) {
	if len(data) == 0 {
		return nil, false
	}

	cMedia := C.getMediaInfo((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))

	videoInfo := &MediaInfo{
		Format:   C.GoString(&cMedia.format[0]),
		Width:    int(cMedia.width),
		Height:   int(cMedia.height),
		Duration: float64(cMedia.duration),
	}

	ok := int(cMedia.ok)
	if ok != 1 {
		return nil, false
	}

	return videoInfo, true
}

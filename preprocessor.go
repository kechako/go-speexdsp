package speexdsp

/*
#cgo pkg-config: speexdsp

#include <speex/speex_preprocess.h>
*/
import "C"
import "unsafe"

type Preprocessor struct {
	state *C.SpeexPreprocessState

	sampleRate int
	frameSize  int
}

func FrameSize(frameSizeMs, sampleRate int) int {
	return frameSizeMs * sampleRate / 1000
}

func NewPreprocessor(sampleRate, frameSize int) *Preprocessor {
	state := C.speex_preprocess_state_init(C.int(frameSize), C.int(sampleRate))

	return &Preprocessor{
		state:      state,
		sampleRate: sampleRate,
		frameSize:  frameSize,
	}
}

func (p *Preprocessor) Close() error {
	C.speex_preprocess_state_destroy(p.state)

	return nil
}

func (p *Preprocessor) SampleRate() int {
	return p.sampleRate
}

func (p *Preprocessor) FrameSize() int {
	return p.frameSize
}

func (p *Preprocessor) Run(buf []int16) {
	pbuf := (*C.spx_int16_t)(&buf[0])
	C.speex_preprocess_run(p.state, pbuf)
}

func (p *Preprocessor) SetEchoCanceller(ec *EchoCanceller) {
	ecState := unsafe.Pointer(nil)
	if ec != nil {
		ecState = unsafe.Pointer(ec.state)
	}
	C.speex_preprocess_ctl(p.state, C.SPEEX_PREPROCESS_SET_ECHO_STATE, ecState)
}

func (p *Preprocessor) EnableDenoise(enable bool) {
	denoise := 0
	if enable {
		denoise = 1
	}

	C.speex_preprocess_ctl(p.state, C.SPEEX_PREPROCESS_SET_DENOISE, unsafe.Pointer(&denoise))
}

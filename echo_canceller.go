package speexdsp

/*
#cgo pkg-config: speexdsp

#include <speex/speex_echo.h>
*/
import "C"
import "unsafe"

type EchoCanceller struct {
	state *C.SpeexEchoState

	micChannels     int
	speekerChannels int
	sampleRate      int
	frameSize       int
	filterLength    int
}

func NewEchoCanceller(micChannels, speekerChannels, sampleRate, frameSizeMs, filterLengthMs int) *EchoCanceller {
	frameSize := frameSizeMs * sampleRate / 1000
	filterLength := filterLengthMs * sampleRate / 1000

	var state *C.SpeexEchoState
	if micChannels == 1 && speekerChannels == 1 {
		state = C.speex_echo_state_init(C.int(frameSize), C.int(filterLength))
	} else {
		state = C.speex_echo_state_init_mc(C.int(frameSize), C.int(filterLength), C.int(micChannels), C.int(speekerChannels))
	}

	C.speex_echo_ctl(state, C.SPEEX_ECHO_SET_SAMPLING_RATE, unsafe.Pointer(&sampleRate))

	return &EchoCanceller{
		state:           state,
		micChannels:     micChannels,
		speekerChannels: speekerChannels,
		sampleRate:      sampleRate,
		frameSize:       frameSize,
		filterLength:    filterLength,
	}
}

func (ec *EchoCanceller) Close() error {
	C.speex_echo_state_destroy(ec.state)

	return nil
}

func (ec *EchoCanceller) MicrophoneChannels() int {
	return ec.micChannels
}

func (ec *EchoCanceller) SpeekerChannels() int {
	return ec.speekerChannels
}

func (ec *EchoCanceller) FrameSize() int {
	return ec.frameSize
}

func (ec *EchoCanceller) FilterLength() int {
	return ec.filterLength
}

func (ec *EchoCanceller) Cancellation(rec, echo, out []int16) {
	prec := (*C.spx_int16_t)(&rec[0])
	pecho := (*C.spx_int16_t)(&echo[0])
	pout := (*C.spx_int16_t)(&out[0])
	C.speex_echo_cancellation(ec.state, prec, pecho, pout)
}

func (ec *EchoCanceller) Capture(rec, out []int16) {
	prec := (*C.spx_int16_t)(&rec[0])
	pout := (*C.spx_int16_t)(&out[0])
	C.speex_echo_capture(ec.state, prec, pout)
}

func (ec *EchoCanceller) Playback(play []int16) {
	pplay := (*C.spx_int16_t)(&play[0])
	C.speex_echo_playback(ec.state, pplay)
}

package speexdsp

/*
#cgo pkg-config: speexdsp

#include <speex/speex_resampler.h>
*/
import "C"
import (
	"errors"
)

const (
	ResamplerQualityMax     = C.SPEEX_RESAMPLER_QUALITY_MAX
	ResamplerQualityMin     = C.SPEEX_RESAMPLER_QUALITY_MIN
	DefaultResamplerQuality = C.SPEEX_RESAMPLER_QUALITY_DEFAULT
	ResamplerQualityVoIP    = C.SPEEX_RESAMPLER_QUALITY_VOIP
	ResamplerQualityDesktop = C.SPEEX_RESAMPLER_QUALITY_DESKTOP
)

var (
	errResamplerAllocFailed = errors.New("memory allocation failed")
	errResamplerBadState    = errors.New("bad resampler state")
	errResamplerInvalidArgs = errors.New("invalid argument")
	errResamplerPtrOverlap  = errors.New("input and output buffers overlap")
	errResamplerUnknown     = errors.New("unknown error. bad error code or strange version mismatch")
)

func errorFromRet(ret C.int) (err error) {
	switch ret {
	case C.RESAMPLER_ERR_SUCCESS:
		// success
	case C.RESAMPLER_ERR_ALLOC_FAILED:
		err = errResamplerAllocFailed
	case C.RESAMPLER_ERR_BAD_STATE:
		err = errResamplerBadState
	case C.RESAMPLER_ERR_INVALID_ARG:
		err = errResamplerInvalidArgs
	case C.RESAMPLER_ERR_PTR_OVERLAP:
		err = errResamplerPtrOverlap
	default:
		err = errResamplerUnknown
	}

	return
}

type Resampler struct {
	state *C.SpeexResamplerState

	channels int
}

func NewResampler(channels, inRate, outRate, quality int) (*Resampler, error) {
	if channels < 1 {
		return nil, errResamplerInvalidArgs
	}
	if inRate < 0 {
		return nil, errResamplerInvalidArgs
	}
	if outRate < 0 {
		return nil, errResamplerInvalidArgs
	}
	if quality < 0 || quality > 10 {
		return nil, errResamplerInvalidArgs
	}

	var ret C.int
	state := C.speex_resampler_init(
		C.spx_uint32_t(channels),
		C.spx_uint32_t(inRate),
		C.spx_uint32_t(outRate),
		C.int(quality),
		&ret)
	if state == nil {
		return nil, errorFromRet(ret)
	}

	return &Resampler{
		state:    state,
		channels: channels,
	}, nil
}

func (r *Resampler) Close() error {
	C.speex_resampler_destroy(r.state)

	return nil
}

func (r *Resampler) ProcessFloat(channelIndex int, in []float32, out []float32) (int, int, error) {
	if channelIndex < 0 || channelIndex >= r.channels {
		return 0, 0, errResamplerInvalidArgs
	}

	pin := (*C.float)(&in[0])
	lenin := C.spx_uint32_t(len(in))
	pout := (*C.float)(&out[0])
	lenout := C.spx_uint32_t(len(out))

	ret := C.speex_resampler_process_float(r.state, C.spx_uint32_t(channelIndex), pin, &lenin, pout, &lenout)
	if err := errorFromRet(ret); err != nil {
		return 0, 0, err
	}

	return int(lenin), int(lenout), nil
}

func (r *Resampler) ProcessInt(channelIndex int, in []int16, out []int16) (int, int, error) {
	if channelIndex < 0 || channelIndex >= r.channels {
		return 0, 0, errResamplerInvalidArgs
	}

	pin := (*C.spx_int16_t)(&in[0])
	lenin := C.spx_uint32_t(len(in))
	pout := (*C.spx_int16_t)(&out[0])
	lenout := C.spx_uint32_t(len(out))

	ret := C.speex_resampler_process_int(r.state, C.spx_uint32_t(channelIndex), pin, &lenin, pout, &lenout)
	if err := errorFromRet(ret); err != nil {
		return 0, 0, err
	}

	return int(lenin), int(lenout), nil
}

func (r *Resampler) ProcessInterleavedFloat(in []float32, out []float32) (int, int, error) {
	pin := (*C.float)(&in[0])
	lenin := C.spx_uint32_t(len(in))
	pout := (*C.float)(&out[0])
	lenout := C.spx_uint32_t(len(out))

	ret := C.speex_resampler_process_interleaved_float(r.state, pin, &lenin, pout, &lenout)
	if err := errorFromRet(ret); err != nil {
		return 0, 0, err
	}

	return int(lenin), int(lenout), nil
}

func (r *Resampler) ProcessInterleavedInt(in []int16, out []int16) (int, int, error) {
	pin := (*C.spx_int16_t)(&in[0])
	lenin := C.spx_uint32_t(len(in))
	pout := (*C.spx_int16_t)(&out[0])
	lenout := C.spx_uint32_t(len(out))

	ret := C.speex_resampler_process_interleaved_int(r.state, pin, &lenin, pout, &lenout)
	if err := errorFromRet(ret); err != nil {
		return 0, 0, err
	}

	return int(lenin), int(lenout), nil
}

func (r *Resampler) SetSampleRate(inRate, outRate int) error {
	ret := C.speex_resampler_set_rate(r.state, C.spx_uint32_t(inRate), C.spx_uint32_t(outRate))
	if err := errorFromRet(ret); err != nil {
		return err
	}

	return nil
}

func (r *Resampler) SampleRate() (int, int) {
	var inRate, outRate C.spx_uint32_t
	C.speex_resampler_get_rate(r.state, &inRate, &outRate)

	return int(inRate), int(outRate)
}

func (r *Resampler) SetQuality(quality int) error {
	ret := C.speex_resampler_set_quality(r.state, C.int(quality))
	if err := errorFromRet(ret); err != nil {
		return err
	}

	return nil
}

func (r *Resampler) Quality() int {
	var quality C.int
	C.speex_resampler_get_quality(r.state, &quality)

	return int(quality)
}

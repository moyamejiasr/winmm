/* For doc goto: https://docs.microsoft.com/es-es/windows/win32/api/mmeapi/ */
package winmm

/*
#cgo windows LDFLAGS: -lwinmm
#include <stdio.h>
#include <Windows.h>
*/
import "C"
import "unsafe"

type WaveFormat C.WAVEFORMATEX
const ChMono = 1
const ChStereo = 2
const SpS08kHz = 8000
const SpS11kHz = 11025
const SpS22kHz = 22050
const SpS44kHz = 44100
const BpS08 = 8
const BpS16 = 16

func NewWaveFormat(ch uint16, sps uint32, bps uint16) WaveFormat {
	var wFmt = WaveFormat{
		wFormatTag: C.WAVE_FORMAT_PCM,
		nChannels: C.WORD(ch),
		nSamplesPerSec: C.DWORD(sps),
		wBitsPerSample: C.WORD(bps),
	}
	wFmt.nBlockAlign = wFmt.nChannels * wFmt.wBitsPerSample / 8
	wFmt.nAvgBytesPerSec = wFmt.nSamplesPerSec * C.DWORD(wFmt.nBlockAlign)
	return wFmt
}

const MMNoError = C.MMSYSERR_NOERROR
const WaveMapper = uint32(C.WAVE_MAPPER)

type WaveHeader C.WAVEHDR

func NewWaveHeader(wFmt WaveFormat, sec uint32) *WaveHeader {
	length := wFmt.nSamplesPerSec * C.DWORD(wFmt.nBlockAlign) * C.DWORD(sec)
	return &WaveHeader{
		lpData: C.LPSTR(C.malloc(C.size_t(length))),
		dwBufferLength: C.DWORD(length),
	}
}

func CopyWaveHeader(header *WaveHeader) *WaveHeader {
	return &WaveHeader{
		lpData: header.lpData,
		dwBufferLength: header.dwBufferLength,
	}
}

func (wh *WaveHeader) ClearBuffer() {
	C.memset(unsafe.Pointer(wh.lpData),
		0, C.size_t(wh.dwBufferLength))
}

func (wh *WaveHeader) GetBufferLength() int {
	return int(wh.dwBufferLength)
}

func (wh *WaveHeader) GetBufferSlice() []byte {
	return C.GoBytes(unsafe.Pointer(wh.lpData), C.int(wh.dwBufferLength))
}
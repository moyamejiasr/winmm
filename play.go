/* For doc goto: https://docs.microsoft.com/es-es/windows/win32/api/mmeapi/ */
package winmm

/*
#cgo amd64 CFLAGS: -DAMD64=1
#cgo windows LDFLAGS: -lwinmm
#include <Windows.h>
void CALLBACK CWaveOutProc(HWAVEOUT, UINT, DWORD_PTR, DWORD_PTR, DWORD_PTR);
*/
import "C"
import (
	"errors"
	"unsafe"
)

const MMOOpen = C.WOM_OPEN
const MMODone = C.WOM_DONE
const MMOClose = C.WOM_CLOSE

type WaveOutDevice C.WAVEOUTCAPS

func EnumOutDevices() ([]WaveOutDevice, error)  {
	var devices []WaveOutDevice
	deviceCount := uint32(C.waveOutGetNumDevs())
	for i := uint32(0); i < deviceCount; i++ {
		var device WaveOutDevice
		result := C.waveOutGetDevCaps(C.UINT_PTR(i),
			C.LPWAVEOUTCAPS(unsafe.Pointer(&device)), C.sizeof_WAVEOUTCAPS)
		if result != MMNoError {
			return devices, getInErrText(result)
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (wid *WaveOutDevice) Id() uint32 {
	return uint32(wid.wPid)
}

func (wid *WaveOutDevice) Name() string {
	return C.GoStringN(C.LPSTR(unsafe.Pointer(&wid.szPname)), C.MAXPNAMELEN)
}

type WaveOut struct {
	waveOut	C.HWAVEOUT
	buffer	*WaveHeader
	Channel	chan *WaveHeader
}

//export waveOutProc
func waveOutProc(instance C.DWORD_PTR, msg C.UINT, param C.DWORD_PTR) {
	var wOut = (*WaveOut)(unsafe.Pointer(uintptr(instance)))
	var wH = (*WaveHeader)(unsafe.Pointer(uintptr(param)))

	switch msg {
	case MMOOpen:
		wOut.Channel = make(chan *WaveHeader)
	case MMODone:
		wOut.Channel <- wH
	case MMOClose:
		close(wOut.Channel)
	}
}

func getOutErrText(result C.MMRESULT) error {
	var length = C.size_t(1024)
	var pszText = C.LPSTR(C.malloc(length))
	defer C.free(unsafe.Pointer(pszText))
	C.waveOutGetErrorText(result, pszText, C.UINT(length))
	return errors.New(C.GoStringN(pszText, C.int(length)))
}

func NewWaveOut() *WaveOut {
	return &WaveOut{}
}

func (wo *WaveOut) Open(dId uint32, wFmt WaveFormat) error {
	wFmtPtr := C.LPCWAVEFORMATEX(unsafe.Pointer(&wFmt))
	wOutProcPtr := C.DWORD_PTR(uintptr(C.CWaveOutProc))
	wOutPtr := C.DWORD_PTR(uintptr(unsafe.Pointer(wo)))

	result := C.waveOutOpen(&wo.waveOut, C.UINT(dId),
		wFmtPtr, wOutProcPtr, wOutPtr, C.CALLBACK_FUNCTION)
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) Close() error {
	result := C.waveOutClose(wo.waveOut)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wo *WaveOut) SetVolume(left, right uint16) error {
	// NOTE: Even tho the docs say 0xFFFF is full-volume
	// it's only left one. We should use 0xFFFFFFFF
	result := C.waveOutSetVolume(wo.waveOut,
		C.DWORD(uint32(left) << 16 | uint32(right)))
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) SetPitch(val uint16) error {
	result := C.waveOutSetPitch(wo.waveOut,
		C.DWORD(val))
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) PrepareBuffer(buffer *WaveHeader) error {
	wo.buffer = buffer
	wHdrPtr := C.LPWAVEHDR(unsafe.Pointer(wo.buffer))

	// Prepare header
	result := C.waveOutPrepareHeader(wo.waveOut,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) UnPrepareBuffer() error {
	wHdrPtr := C.LPWAVEHDR(unsafe.Pointer(wo.buffer))

	result := C.waveOutUnprepareHeader(wo.waveOut,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) Write() error {
	wHdrPtr := C.LPWAVEHDR(unsafe.Pointer(wo.buffer))

	result := C.waveOutWrite(wo.waveOut,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return getOutErrText(result)
	}
	return nil
}

func (wo *WaveOut) Pause() error {
	result := C.waveOutPause(wo.waveOut)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wo *WaveOut) Restart() error {
	result := C.waveOutRestart(wo.waveOut)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wo *WaveOut) Reset() error {
	result := C.waveOutReset(wo.waveOut)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}
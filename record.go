/* For doc goto: https://docs.microsoft.com/es-es/windows/win32/api/mmeapi/ */
package winmm

/*
#cgo amd64 CFLAGS: -DAMD64=1
#cgo windows LDFLAGS: -lwinmm
#include <Windows.h>
void CALLBACK CWaveInProc(HWAVEIN, UINT, DWORD_PTR, DWORD_PTR, DWORD_PTR);
*/
import "C"
import (
	"errors"
	"unsafe"
)

const MMIOpen = C.WIM_OPEN
const MMIData = C.WIM_DATA
const MMIClose = C.WIM_CLOSE

type WaveInDevice C.WAVEINCAPS

func EnumInDevices() ([]WaveInDevice, error)  {
	var devices []WaveInDevice
	deviceCount := uint32(C.waveInGetNumDevs())
	for i := uint32(0); i < deviceCount; i++ {
		var device WaveInDevice
		result := C.waveInGetDevCaps(C.UINT_PTR(i),
			C.LPWAVEINCAPS(unsafe.Pointer(&device)), C.sizeof_WAVEINCAPS)
		if result != MMNoError {
			return devices, getInErrText(result)
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (wid *WaveInDevice) Id() uint32 {
	return uint32(wid.wPid)
}

func (wid *WaveInDevice) Name() string {
	return C.GoStringN(C.LPSTR(unsafe.Pointer(&wid.szPname)), C.MAXPNAMELEN)
}

type WaveIn struct {
	waveIn  C.HWAVEIN
	buffers []*WaveHeader
	Channel chan *WaveHeader
}

//export waveInProc
func waveInProc(instance C.DWORD_PTR, msg C.UINT, param C.DWORD_PTR) {
	var wIn = (*WaveIn)(unsafe.Pointer(uintptr(instance)))
	var wH = (*WaveHeader)(unsafe.Pointer(uintptr(param)))

	switch msg {
	case MMIOpen:
		wIn.Channel = make(chan *WaveHeader)
	case MMIData:
		wIn.Channel <- wH
	case MMIClose:
		close(wIn.Channel)
	}
}

func getInErrText(result C.MMRESULT) error {
	var length = C.size_t(1024)
	var pszText = C.LPSTR(C.malloc(length))
	defer C.free(unsafe.Pointer(pszText))
	C.waveInGetErrorText(result, pszText, C.UINT(length))
	return errors.New(C.GoStringN(pszText, C.int(length)))
}

func NewWaveIn() *WaveIn {
	return &WaveIn{}
}

func (wi *WaveIn) Open(dId uint32, wFmt WaveFormat) error {
	wFmtPtr := C.LPCWAVEFORMATEX(unsafe.Pointer(&wFmt))
	wInProcPtr := C.DWORD_PTR(uintptr(C.CWaveInProc))
	wInPtr := C.DWORD_PTR(uintptr(unsafe.Pointer(wi)))

	result := C.waveInOpen(&wi.waveIn, C.UINT(dId),
		wFmtPtr, wInProcPtr, wInPtr, C.CALLBACK_FUNCTION)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) Close() error {
	result := C.waveInClose(wi.waveIn)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) PrepareBuffer(buffer *WaveHeader) (int, error) {
	var key = len(wi.buffers)
	wi.buffers = append(wi.buffers, buffer)
	wHdrPtr := C.LPWAVEHDR(
		unsafe.Pointer(buffer))

	// Prepare header
	result := C.waveInPrepareHeader(wi.waveIn,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return -1, getInErrText(result)
	}
	return key, nil
}

func (wi *WaveIn) UnPrepareBuffer(i int) error {
	wHdrPtr := C.LPWAVEHDR(
		unsafe.Pointer(wi.buffers[i]))

	result := C.waveInUnprepareHeader(wi.waveIn,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) AddBuffer(i int) error {
	wHdrPtr := C.LPWAVEHDR(
		unsafe.Pointer(wi.buffers[i]))

	// Append the buffers to the recording
	result := C.waveInAddBuffer(wi.waveIn,
		wHdrPtr, C.sizeof_WAVEHDR)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) GetBuffer(i int) *WaveHeader {
	return wi.buffers[i]
}

func (wi *WaveIn) RemoveBuffer(i int) {
	wi.buffers = append(wi.buffers[:i], wi.buffers[i+1:]...)
}

func (wi *WaveIn) Start() error {
	result := C.waveInStart(wi.waveIn)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) Stop() error {
	result := C.waveInStop(wi.waveIn)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}

func (wi *WaveIn) Reset() error {
	result := C.waveInReset(wi.waveIn)
	if result != MMNoError {
		return getInErrText(result)
	}
	return nil
}
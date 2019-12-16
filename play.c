#include "_cgo_export.h"
void CALLBACK CWaveOutProc(HWAVEOUT h, UINT m, DWORD_PTR i, DWORD_PTR p1, DWORD_PTR p2)
{
	waveOutProc(i, m, p1);
}
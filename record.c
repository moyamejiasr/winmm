#include "_cgo_export.h"
void CALLBACK CWaveInProc(HWAVEIN h, UINT m, DWORD_PTR i, DWORD_PTR p1, DWORD_PTR p2)
{
	waveInProc(i, m, p1);
}
package main

import (
	"fmt"
	"github.com/Onelio/winmm"
	"time"
)

func main() {
	/* IMPORTANT: In this example I ignore all of the
	errors for the sake of simplicity but you should not.
	*/

	// List Input Devices
	iList, _ := winmm.EnumInDevices()
	for _, item := range iList {
		fmt.Println(item.Id(), item.Name())
	}
	// List Output Devices
	oList, _ := winmm.EnumOutDevices()
	for _, item := range oList {
		fmt.Println(item.Id(), item.Name())
	}
	_, _ = fmt.Scanln()

	// 1- Create WaveFormat for recording
	wFmt := winmm.NewWaveFormat(winmm.ChStereo, winmm.SpS44kHz, winmm.BpS16)

	// Create WaveHeader that contains the buffer
	// for reading with specified format
	// and length of one second
	rBuffer := winmm.NewWaveHeader(wFmt, 2)
	// Copy raw buffer from reading in a new
	// WaveHeader struct for playing
	pBuffer := winmm.CopyWaveHeader(rBuffer)

	// Create the playing struct and open it
	play := winmm.NewWaveOut()
	_ = play.Open(winmm.WaveMapper, wFmt)
	_ = play.SetVolume(0xFFFF, 0xFFFF)
	// Prepare the buffer for playing
	_ = play.PrepareBuffer(pBuffer)

	// Create the recording struct and open it
	record := winmm.NewWaveIn()
	_ = record.Open(winmm.WaveMapper, wFmt)
	// Prepare the buffer for recording and since
	// recording can have multiple buffers we get
	// its key for using later on
	key, _ := record.PrepareBuffer(rBuffer)

	fmt.Println("Recording... ")
	// Add the buffer by it's key to the
	// recording device and start recording
	_ = record.AddBuffer(key)
	_ = record.Start()
	// Give a delay for the playing
	time.Sleep(200 * time.Millisecond)
	// Start playing
	_ = play.Write()

	// Scan for enter key to close
	var shouldClose = false
	go func() {
		_, _ = fmt.Scanln()
		shouldClose = true
	}()

	for !shouldClose {
		select {
			// Restart when one second recorded/played
			case _ = <- record.Channel:
				_ = record.AddBuffer(key)
				_ = record.Start()
				break
			case _ = <- play.Channel:
				_ = play.Write()
				break
		}
	}
	// Un-prepare all at the end
	_ = record.UnPrepareBuffer(key)
	_ = play.UnPrepareBuffer()
	// Close the devices
	_ = record.Close()
	_ = play.Close()
}

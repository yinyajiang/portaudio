package simpleutil

import (
	"github.com/yinyajiang/portaudio"
)

//GetDefaultInputDeviceName ...
func GetDefaultInputDeviceName() (name string, err error) {
	err = portaudio.Initialize()
	if err != nil {
		return
	}
	defer portaudio.Terminate()

	defdev, err := portaudio.DefaultInputDevice()
	if err != nil {
		return
	}
	name = defdev.Name
	return
}

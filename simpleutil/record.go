package simpleutil

import (
	"context"
	"encoding/binary"

	"io"

	"github.com/yinyajiang/portaudio"
)

//RecordAiff ...
func RecordAiff(ctx context.Context, w io.WriteSeeker) (err error) {
	err = writeHAiffeader(w)
	if err != nil {
		return
	}
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = w.Seek(4, 0)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, int32(totalBytes))
		_, err = w.Seek(22, 0)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, int32(nSamples))
		_, err = w.Seek(42, 0)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, int32(4*nSamples+8))
	}()

	err = portaudio.Initialize()
	if err != nil {
		return
	}

	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		return
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return
	}

loop:
	for {
		err = stream.Read()
		if err != nil {
			err = nil
			return
		}
		err = binary.Write(w, binary.BigEndian, in)
		if err != nil {
			return
		}
		nSamples += len(in)
		select {
		case <-ctx.Done():
			break loop
		default:
		}
	}
	stream.Stop()
	return
}

func writeHAiffeader(w io.WriteSeeker) (err error) {
	// form chunk
	_, err = w.Write([]byte("FORM"))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(0)) //total bytes
	if err != nil {
		return
	}
	_, err = w.Write([]byte("AIFF"))
	if err != nil {
		return
	}

	// common chunk
	_, err = w.Write([]byte("COMM"))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(18)) //size
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(1)) //channels
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(0)) //number of samples
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(32)) //bits per sample
	if err != nil {
		return
	}
	_, err = w.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	if err != nil {
		return
	}

	// sound chunk
	_, err = w.Write([]byte("SSND"))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(0)) //size
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(0)) //offset
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(0)) //block
	return
}

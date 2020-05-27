package gumble

import (
	"time"
)

const (
	// "AudioMaximumSampleRate" is the maximum audio sample rate (in Hertz) for
	// incoming and outgoing audio.
	AudioMaximumSampleRate = 48000
	// "AudioSampleRate" is the audio sample rate (in Hertz) for incoming and
	// outgoing audio. Supported rates are 48000, 24000, 16000, 12000 and 8000 Hz.
	// Higher audio samping rate produce higher quality, but results with
	// higher CPU usage. Always make a compromise between audio quality,
	// and CPU usage for your particular application and hardware type.
	AudioSampleRate = 48000
	// "AudioDefaultIntervalMS is" the default interval in milliseconds that audio
	// packets are sent at. Time intervals from 10 to 120 ms can be used.
	// (e.g. 10, 20, 40, 60, 80, 100, 120). Shorter interval result with less
	// latency but require a higher network bandwidth. Always make a compromise
	// between latency and network bandwidth.
	AudioDefaultIntervalMS = 10
	// "AudioDefaultInterval" is the default interval that audio packets are sent
	// at.
	AudioDefaultInterval = AudioDefaultIntervalMS * time.Millisecond
	// "AudioDefaultFrameSize" is the number of audio frames that should be sent in
	// an AudioDefaultInterval window.
	AudioDefaultFrameSize = (AudioSampleRate * AudioDefaultIntervalMS) / 1000
	// "AudioMaximumFrameSize" is the maximum audio frame size from another user
	// that will be processed.
	AudioMaximumFrameSize = AudioMaximumSampleRate / 1000 * 120
	// "AudioDefaultDataBytes" is the default number of bytes that an audio frame
	// can use.
	AudioDefaultDataBytes = AudioSampleRate * AudioDefaultIntervalMS / 1000 * 8
	// AudioDefaultDataBytes = 120
	// "AudioChannels" is the number of audio channels that are contained in an
	// audio stream.
	AudioChannels = 1
)

// AudioListener is the interface that must be implemented by types wishing to
// receive incoming audio data from the server.
//
// OnAudioStream is called when an audio stream for a user starts. It is the
// implementer's responsibility to continuously process AudioStreamEvent.C
// until it is closed.

type AudioListener interface {
	OnAudioStream(e *AudioStreamEvent)
}

// AudioStreamEvent is event that is passed to AudioListener.OnAudioStream.

type AudioStreamEvent struct {
	Client *Client
	User   *User
	C      <-chan *AudioPacket
}

// AudioBuffer is a slice of PCM audio samples.

type AudioBuffer []int16

func (a AudioBuffer) writeAudio(client *Client, seq int64, final bool) error {
	encoder := client.AudioEncoder
	if encoder == nil {
		return nil
	}
	dataBytes := client.Config.AudioDataBytes
	raw, err := encoder.Encode(a, len(a), dataBytes)
	if final {
		defer encoder.Reset()
	}
	if err != nil {
		return err
	}
	var targetID byte
	if target := client.VoiceTarget; target != nil {
		targetID = byte(target.ID)
	}
	// TODO: re-enable positional audio
	return client.Conn.WriteAudio(byte(4), targetID, seq, final, raw, nil, nil, nil)
}

// AudioPacket contains incoming audio samples and information.

type AudioPacket struct {
	Client *Client
	Sender *User
	Target *VoiceTarget
	AudioBuffer
	HasPosition bool
	X, Y, Z     float32
}

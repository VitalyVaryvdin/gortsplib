package format

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pion/rtp"

	"github.com/bluenviron/gortsplib/v4/pkg/format/rtpmpeg4video"
)

// MPEG4Video is an alias for MPEG4VideoES.
type MPEG4Video = MPEG4VideoES

// MPEG4VideoES is a RTP format for a MPEG-4 Video codec.
// Specification: https://datatracker.ietf.org/doc/html/rfc6416#section-7.1
type MPEG4VideoES struct {
	PayloadTyp     uint8
	ProfileLevelID int
	Config         []byte
}

func (f *MPEG4VideoES) unmarshal(
	payloadType uint8, _ string, _ string,
	_ string, fmtp map[string]string,
) error {
	f.PayloadTyp = payloadType
	f.ProfileLevelID = 1 // default value defined by specification

	for key, val := range fmtp {
		switch key {
		case "profile-level-id":
			tmp, err := strconv.ParseUint(val, 10, 31)
			if err != nil {
				return fmt.Errorf("invalid profile-level-id: %v", val)
			}

			f.ProfileLevelID = int(tmp)

		case "config":
			var err error
			f.Config, err = hex.DecodeString(val)
			if err != nil {
				return fmt.Errorf("invalid config: %v", val)
			}
		}
	}

	return nil
}

// Codec implements Format.
func (f *MPEG4VideoES) Codec() string {
	return "MPEG-4 Video"
}

// ClockRate implements Format.
func (f *MPEG4VideoES) ClockRate() int {
	return 90000
}

// PayloadType implements Format.
func (f *MPEG4VideoES) PayloadType() uint8 {
	return f.PayloadTyp
}

// RTPMap implements Format.
func (f *MPEG4VideoES) RTPMap() string {
	return "MP4V-ES/90000"
}

// FMTP implements Format.
func (f *MPEG4VideoES) FMTP() map[string]string {
	fmtp := map[string]string{
		"profile-level-id": strconv.FormatInt(int64(f.ProfileLevelID), 10),
		"config":           strings.ToUpper(hex.EncodeToString(f.Config)),
	}

	return fmtp
}

// PTSEqualsDTS implements Format.
func (f *MPEG4VideoES) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}

// CreateDecoder creates a decoder able to decode the content of the format.
func (f *MPEG4VideoES) CreateDecoder() (*rtpmpeg4video.Decoder, error) {
	d := &rtpmpeg4video.Decoder{}

	err := d.Init()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// CreateEncoder creates an encoder able to encode the content of the format.
func (f *MPEG4VideoES) CreateEncoder() (*rtpmpeg4video.Encoder, error) {
	e := &rtpmpeg4video.Encoder{
		PayloadType: f.PayloadTyp,
	}

	err := e.Init()
	if err != nil {
		return nil, err
	}

	return e, nil
}
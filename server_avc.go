package gortmp

import "errors"

func parseAVCConfig(data []byte, s *Stream) error {
	if len(data) < 7 {
		return errors.New("avc config too short")
	}

	// skip: version(1), profile(1), compat(1), level(1), lengthSizeMinusOne(1)
	i := 5

	numSPS := int(data[i] & 0x1f)
	i++

	for j := 0; j < numSPS; j++ {
		l := int(data[i])<<8 | int(data[i+1])
		i += 2
		s.SPS = data[i : i+l]
		i += l
	}

	numPPS := int(data[i])
	i++

	for j := 0; j < numPPS; j++ {
		l := int(data[i])<<8 | int(data[i+1])
		i += 2
		s.PPS = data[i : i+l]
		i += l
	}

	return nil
}

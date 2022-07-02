package ethernet

import (
	"encoding/binary"
	"io"
)

type MPLS struct {
	Label         uint32
	TrafficClass  uint8
	BottomOfStack bool
	TTL           uint8
}

// MarshalBinary allocates a byte slice and marshals a MPLS into binary form.
func (v *MPLS) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	_, err := v.read(b)
	return b, err
}

// read reads data from a MPLS into b.  read is used to marshal a MPLS into
// binary form, but does not allocate on its own.
func (v *MPLS) read(b []byte) (int, error) {
	// 3 bits: priority
	ub := v.Label << 12
	ub |= uint32(v.TrafficClass) << 9
	var bottom uint32
	if v.BottomOfStack {
		bottom = 1
	}
	ub |= bottom << 8
	ub |= uint32(v.TTL)

	binary.BigEndian.PutUint32(b, ub)
	return 4, nil
}

// UnmarshalBinary unmarshals a byte slice into a MPLS.
func (v *MPLS) UnmarshalBinary(b []byte) error {
	if len(b) != 4 {
		return io.ErrUnexpectedEOF
	}

	ub := binary.BigEndian.Uint32(b[0:4])
	v.Label = ub >> 12
	v.TrafficClass = uint8((ub >> 9) & 0b111)
	v.BottomOfStack = ub&0x100 != 0
	v.TTL = uint8(ub & 0xFF)

	return nil
}

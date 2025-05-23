package packets

import (
	"io"
)

// PingResp is an internal representation of the fields of the PINGRESP MQTT packet.
type PingResp struct {
	FixedHeader
}

func (pkt *PingResp) String() string {
	return pkt.FixedHeader.String()
}

func (pkt *PingResp) Pack(w io.Writer) error {
	bytes := pkt.FixedHeader.Encode()
	_, err := w.Write(bytes)

	return err
}

func (pkt *PingResp) Unpack(r io.Reader, v byte) error {
	return nil
}

func (pkt *PingResp) Details() Details {
	return Details{Type: PingRespType, ID: 0, Qos: 0}
}

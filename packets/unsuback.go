package packets

import (
	"fmt"
	"io"

	codec "github.com/dborovcanin/mqtt/packets/codec"
)

// UnSubAck is an internal representation of the fields of the UNSUBACK MQTT packet.

type UnSubAck struct {
	FixedHeader
	// Variable Header
	ID         uint16
	Properties *BasicProperties
	// Payload
	ReasonCodes *[]byte
}

func (pkt *UnSubAck) String() string {
	return fmt.Sprintf("%s\npacket_id: %d\n", pkt.FixedHeader, pkt.ID)
}

func (pkt *UnSubAck) Pack(w io.Writer) error {
	bytes := codec.EncodeUint16(pkt.ID)
	if pkt.Properties != nil {
		props := pkt.Properties.Encode()
		l := len(props)
		proplen := codec.EncodeVBI(l)
		bytes = append(bytes, proplen...)
		if l > 0 {
			bytes = append(bytes, props...)
		}
	}
	if pkt.ReasonCodes != nil {
		bytes = append(bytes, *pkt.ReasonCodes...)
	}
	// Take care size is calculated properly if someone tempered with the packet.
	pkt.FixedHeader.RemainingLength = len(bytes)
	bytes = append(pkt.FixedHeader.Encode(), bytes...)
	_, err := w.Write(bytes)

	return err
}

func (pkt *UnSubAck) Unpack(r io.Reader, v byte) error {
	var err error
	pkt.ID, err = codec.DecodeUint16(r)
	if err != nil {
		return err
	}
	if v == V5 {
		p := BasicProperties{}
		length, err := codec.DecodeVBI(r)
		if err != nil {
			return err
		}
		if length != 0 {
			if err := p.Unpack(r); err != nil {
				return err
			}
			pkt.Properties = &p
		}
		rc, err := codec.DecodeBytes(r)
		if err != nil {
			return err
		}
		pkt.ReasonCodes = &rc
	}

	return nil
}

func (pkt *UnSubAck) Details() Details {
	return Details{Type: SubAckType, ID: pkt.ID, Qos: 0}
}

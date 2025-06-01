package packets

import (
	"bytes"
	"fmt"
	"io"

	codec "github.com/dborovcanin/mqtt/packets/codec"
)

// Disconnect is an internal representation of the fields of the DISCONNECT MQTT packet
type Disconnect struct {
	FixedHeader
	// Variable Header
	Properties *DisconnectProperties
}

type DisconnectProperties struct {
	// SessionExpiryInterval is the time in seconds after a client disconnects
	// that the server should retain the session information (subscriptions etc).
	SessionExpiryInterval *uint32
	// ReasonString is a UTF8 string representing the reason associated with
	// this response, intended to be human readable for diagnostic purposes.
	ReasonString string
	// User is a slice of user provided properties (key and value).
	User []User
	// ServerReference is a UTF8 string indicating another server the client
	// can use.
	ServerReference string
}

func (p *DisconnectProperties) Unpack(r io.Reader) error {
	length, err := codec.DecodeVBI(r)
	if err != nil {
		return err
	}
	if length == 0 {
		return nil
	}
	for {
		prop, err := codec.DecodeByte(r)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch prop {
		case SessionExpiryIntervalProp:
			sei, err := codec.DecodeUint32(r)
			if err != nil {
				return err
			}
			p.SessionExpiryInterval = &sei
		case ReasonStringProp:
			p.ReasonString, err = codec.DecodeString(r)
			if err != nil {
				return err
			}
		case UserProp:
			k, err := codec.DecodeString(r)
			if err != nil {
				return err
			}
			v, err := codec.DecodeString(r)
			if err != nil {
				return err
			}
			p.User = append(p.User, User{k, v})
		case ServerReferenceProp:
			p.ServerReference, err = codec.DecodeString(r)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid property type %d for disconnect packet", prop)
		}
	}
}

func (p *DisconnectProperties) Encode() []byte {
	var ret []byte
	if p.SessionExpiryInterval != nil {
		ret = append(ret, codec.EncodeUint32(*p.SessionExpiryInterval)...)
	}
	if p.ReasonString != "" {
		ret = append(ret, codec.EncodeBytes([]byte(p.ReasonString))...)
	}
	if len(p.User) > 0 {
		for _, u := range p.User {
			ret = append(ret, codec.EncodeBytes([]byte(u.Key))...)
			ret = append(ret, codec.EncodeBytes([]byte(u.Value))...)
		}
	}
	if p.ServerReference != "" {
		ret = append(ret, codec.EncodeBytes([]byte(p.ServerReference))...)
	}

	return ret
}

func (pkt *Disconnect) String() string {
	return pkt.FixedHeader.String()
}

func (pkt *Disconnect) Encode() []byte {
	ret := []byte{}
	if pkt.Properties != nil {
		if pkt.Properties != nil {
			props := pkt.Properties.Encode()
			l := len(props)
			proplen := codec.EncodeVBI(l)
			ret = append(ret, proplen...)
			if l > 0 {
				ret = append(ret, props...)
			}
		}
	}
	// Take care size is calculated properly if someone tempered with the packet.
	pkt.FixedHeader.RemainingLength = len(ret)
	ret = append(pkt.FixedHeader.Encode(), ret...)

	return ret
}

func (pkt *Disconnect) Pack(w io.Writer) error {
	_, err := w.Write(pkt.Encode())
	return err
}

func (pkt *Disconnect) Unpack(r io.Reader, v byte) error {
	if v == V5 {
		length, err := codec.DecodeVBI(r)
		if err != nil {
			return err
		}
		if length == 0 {
			return nil
		}
		buf := make([]byte, length)
		if _, err := r.Read(buf); err != nil {
			return err
		}
		p := DisconnectProperties{}
		props := bytes.NewReader(buf)
		if err := p.Unpack(props); err != nil {
			return err
		}
		pkt.Properties = &p
	}
	return nil
}

func (pkt *Disconnect) Details() Details {
	return Details{Type: DisconnectType, ID: 0, Qos: 0}
}

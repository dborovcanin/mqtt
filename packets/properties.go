package packets

import (
	"fmt"
	"io"

	"github.com/dborovcanin/mqtt/packets/codec"
)

// PropPayloadFormat, etc are the list of property codes for the
// MQTT packet properties
const (
	PayloadFormatProp          byte = 1
	MessageExpiryProp          byte = 2
	ContentTypeProp            byte = 3
	ResponseTopicProp          byte = 8
	CorrelationDataProp        byte = 9
	SubscriptionIdentifierProp byte = 11
	SessionExpiryIntervalProp  byte = 17
	AssignedClientIDProp       byte = 18
	ServerKeepAliveProp        byte = 19
	AuthMethodProp             byte = 21
	AuthDataProp               byte = 22
	RequestProblemInfoProp     byte = 23
	WillDelayIntervalProp      byte = 24
	RequestResponseInfoProp    byte = 25
	ResponseInfoProp           byte = 26
	ServerReferenceProp        byte = 28
	ReasonStringProp           byte = 31
	ReceiveMaximumProp         byte = 33
	TopicAliasMaximumProp      byte = 34
	TopicAliasProp             byte = 35
	MaximumQOSProp             byte = 36
	RetainAvailableProp        byte = 37
	UserProp                   byte = 38
	MaximumPacketSizeProp      byte = 39
	WildcardSubAvailableProp   byte = 40
	SubIDAvailableProp         byte = 41
	SharedSubAvailableProp     byte = 42
)

type User struct {
	Key, Value string
}

type BasicProperties struct {
	// ReasonString is a UTF8 string representing the reason associated with
	// this response, intended to be human readable for diagnostic purposes.
	ReasonString string
	// User is a slice of user provided properties (key and value).
	User []User
}

func (p *BasicProperties) Unpack(r io.Reader) error {
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
		default:
			return fmt.Errorf("invalid property type %d", prop)
		}
	}
}

func (p *BasicProperties) Encode() []byte {
	var ret []byte
	if p.ReasonString != "" {
		ret = append(ret, codec.EncodeBytes([]byte(p.ReasonString))...)
	}
	if len(p.User) > 0 {
		for _, u := range p.User {
			ret = append(ret, codec.EncodeBytes([]byte(u.Key))...)
			ret = append(ret, codec.EncodeBytes([]byte(u.Value))...)
		}
	}

	return ret
}

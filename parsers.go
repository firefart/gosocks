package socks

import (
	"encoding/binary"
	"fmt"
)

/*
	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+

o  VER    protocol version: X'05'
o  CMD

	o  CONNECT X'01'
	o  BIND X'02'
	o  UDP ASSOCIATE X'03'

o  RSV    RESERVED
o  ATYP   address type of following address

	o  IP V4 address: X'01'
	o  DOMAINNAME: X'03'
	o  IP V6 address: X'04'

o  DST.ADDR       desired destination address
o  DST.PORT desired destination port in network octet

	order

In an address field (DST.ADDR, BND.ADDR), the ATYP field specifies
the type of address contained within the field:

	o  X'01'

the address is a version-4 IP address, with a length of 4 octets

	o  X'03'

the address field contains a fully-qualified domain name.  The first
octet of the address field contains the number of octets of name that
follow, there is no terminating NUL octet.

	o  X'04'

the address is a version-6 IP address, with a length of 16 octets.
*/
func parseRequest(buf []byte) (*Request, *Error) {
	r := &Request{}
	if len(buf) < 7 {
		return nil, &Error{Reason: RequestReplyConnectionRefused, Err: fmt.Errorf("invalid request header length (%d)", len(buf))}
	}
	version := buf[0]
	switch version {
	case byte(Version4):
		r.Version = Version4
	case byte(Version5):
		r.Version = Version5
	default:
		return nil, &Error{Reason: RequestReplyConnectionRefused, Err: fmt.Errorf("Invalid Socks version %#x", version)}
	}
	cmd := buf[1]
	switch cmd {
	case byte(RequestCmdConnect):
		r.Command = RequestCmdConnect
	// case byte(RequestCmdBind):
	// 	r.Command = RequestCmdBind
	// case byte(RequestCmdAssociate):
	// 	r.Command = RequestCmdAssociate
	default:
		return nil, &Error{Reason: RequestReplyCommandNotSupported, Err: fmt.Errorf("Command %#x not supported", cmd)}
	}
	addresstype := buf[3]
	switch addresstype {
	case byte(RequestAddressTypeIPv4):
		r.AddressType = RequestAddressTypeIPv4
	case byte(RequestAddressTypeIPv6):
		r.AddressType = RequestAddressTypeIPv6
	case byte(RequestAddressTypeDomainname):
		r.AddressType = RequestAddressTypeDomainname
	default:
		return nil, &Error{Reason: RequestReplyAddressTypeNotSupported, Err: fmt.Errorf("AddressType %#x not supported", addresstype)}
	}

	switch r.AddressType {
	case RequestAddressTypeIPv4:
		r.DestinationAddress = buf[4:8]
		p := buf[8:10]
		r.DestinationPort = binary.BigEndian.Uint16(p)
	case RequestAddressTypeIPv6:
		r.DestinationAddress = buf[4:20]
		p := buf[20:22]
		r.DestinationPort = binary.BigEndian.Uint16(p)
	case RequestAddressTypeDomainname:
		addrLen := buf[4]
		r.DestinationAddress = buf[5 : 5+addrLen]
		p := buf[5+addrLen : 5+addrLen+2]
		r.DestinationPort = binary.BigEndian.Uint16(p)
	default:
		return nil, &Error{Reason: RequestReplyAddressTypeNotSupported, Err: fmt.Errorf("AddressType %#x not supported", addresstype)}
	}

	return r, nil
}

func parseHeader(buf []byte) (Header, error) {
	h := Header{}
	if len(buf) < 3 {
		return h, fmt.Errorf("invalid socks header")
	}
	version := buf[0]
	switch version {
	case byte(Version4):
		h.Version = Version4
	case byte(Version5):
		h.Version = Version5
	default:
		return h, fmt.Errorf("could not get socks version from header")
	}
	numMethods := buf[1]
	if len(buf) < int(numMethods)+2 {
		return h, fmt.Errorf("invalid socks header")
	}
	h.Methods = make([]byte, numMethods)
	for i := 0; i < int(numMethods); i++ {
		meth := buf[2+i]
		h.Methods = append(h.Methods, meth)
	}
	return h, nil
}

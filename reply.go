package socks

import (
	"bytes"
	"encoding/binary"
)

/*
	        +----+-----+-------+------+----------+----------+
	        |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	        +----+-----+-------+------+----------+----------+
	        | 1  |  1  | X'00' |  1   | Variable |    2     |
	        +----+-----+-------+------+----------+----------+

	     Where:

	          o  VER    protocol version: X'05'
	          o  REP    Reply field:
	             o  X'00' succeeded
	             o  X'01' general SOCKS server failure
	             o  X'02' connection not allowed by ruleset
	             o  X'03' Network unreachable
	             o  X'04' Host unreachable
	             o  X'05' Connection refused
	             o  X'06' TTL expired
	             o  X'07' Command not supported
	             o  X'08' Address type not supported
	             o  X'09' to X'FF' unassigned
	          o  RSV    RESERVED
	          o  ATYP   address type of following address
	             o  IP V4 address: X'01'
	             o  DOMAINNAME: X'03'
	             o  IP V6 address: X'04'
	          o  BND.ADDR       server bound address
	          o  BND.PORT       server bound port in network octet order

		 Fields marked RESERVED (RSV) must be set to X'00'.
		 CONNECT

	   In the reply to a CONNECT, BND.PORT contains the port number that the
	   server assigned to connect to the target host, while BND.ADDR
	   contains the associated IP address.  The supplied BND.ADDR is often
	   different from the IP address that the client uses to reach the SOCKS
	   server, since such servers are often multi-homed.  It is expected
	   that the SOCKS server will use DST.ADDR and DST.PORT, and the
	   client-side source address and port in evaluating the CONNECT
	   request.
*/
func requestReply(request *Request, reply RequestReplyReason) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := binary.Write(buffer, binary.BigEndian, Version5.Value()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, reply.Value()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, byte(0x00)); err != nil {
		return nil, err
	}

	if request != nil {
		if err := binary.Write(buffer, binary.BigEndian, request.AddressType.Value()); err != nil {
			return nil, err
		}

		// type
		switch request.AddressType {
		case RequestAddressTypeIPv4:
			if err := binary.Write(buffer, binary.BigEndian, request.DestinationAddress); err != nil {
				return nil, err
			}
		case RequestAddressTypeIPv6:
			if err := binary.Write(buffer, binary.BigEndian, request.DestinationAddress); err != nil {
				return nil, err
			}
		default:
			if err := binary.Write(buffer, binary.BigEndian, byte(len(request.DestinationAddress))); err != nil {
				return nil, err
			}
			if err := binary.Write(buffer, binary.BigEndian, request.DestinationAddress); err != nil {
				return nil, err
			}
		}

		if err := binary.Write(buffer, binary.BigEndian, request.DestinationPort); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Write(buffer, binary.BigEndian, RequestAddressTypeIPv4.Value()); err != nil {
			return nil, err
		}
		if err := binary.Write(buffer, binary.BigEndian, []byte{0, 0, 0, 0}); err != nil {
			return nil, err
		}
		if err := binary.Write(buffer, binary.BigEndian, uint16(0)); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

package socks

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"strconv"
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
func requestReply(in net.Addr, reply RequestReplyReason) ([]byte, error) {
	var buf []byte
	buf = append(buf, Version5.Value())
	buf = append(buf, reply.Value())
	// reserved
	buf = append(buf, 0x00)

	if in != nil {
		host, port, err := net.SplitHostPort(in.String())
		if err != nil {
			return nil, err
		}
		ip, err := netip.ParseAddr(host)
		if err != nil {
			return nil, err
		}

		// type
		if ip.Is4() {
			buf = append(buf, RequestAddressTypeIPv4.Value())
		} else if ip.Is6() {
			buf = append(buf, RequestAddressTypeIPv6.Value())
		} else {
			return nil, fmt.Errorf("ip %s invalid", ip.String())
		}

		buf = append(buf, ip.AsSlice()...)
		portInt, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			return nil, err
		}
		var portByte = make([]byte, 2)
		binary.BigEndian.PutUint16(portByte, uint16(portInt))
		buf = append(buf, portByte...)
	} else {
		// type
		buf = append(buf, RequestAddressTypeIPv4.Value())
		// error reply
		buf = append(buf, []byte{0, 0, 0, 0}...)
	}
	return buf, nil
}

package socks

import (
	"fmt"
	"net"
)

// Header holds a Socks5 header
type Header struct {
	Version Version
	Methods []byte
}

// Request holds a Socks5 request
type Request struct {
	Version            Version
	Command            RequestCmd
	Reserved           byte
	AddressType        RequestAddressType
	DestinationAddress []byte
	DestinationPort    uint16
}

func (r Request) getDestinationString() string {
	switch r.AddressType {
	case RequestAddressTypeDomainname:
		return fmt.Sprintf("%s:%d", r.DestinationAddress, r.DestinationPort)
	case RequestAddressTypeIPv4:
		ip := net.IP(r.DestinationAddress)
		return fmt.Sprintf("%s:%d", ip.String(), r.DestinationPort)
	case RequestAddressTypeIPv6:
		ip := net.IP(r.DestinationAddress)
		return fmt.Sprintf("%s:%d", ip.String(), r.DestinationPort)
	default:
		return fmt.Sprintf("Address type %d not implemented", r.AddressType)
	}
}

// Methods holds the socks5 msethod
type Methods uint8

const (
	// MethodNoAuthRequired means the socks proxy requires no auth
	MethodNoAuthRequired = 0x00
	// MethodGSSAPI means the socks proxy requires authentication with GSSAPI
	MethodGSSAPI = 0x01
	// MethodUsernamePassword means the socks proxy requires authentication with username and passowrd
	MethodUsernamePassword = 0x02
	// MethodNoAcceptableMethods means the socks proxy does not implement any of the requested methods
	MethodNoAcceptableMethods = 0xff
)

// Version holds the socks5 version
type Version uint8

// Value gets the real value of the Version
func (v Version) Value() uint8 {
	return uint8(v)
}

const (
	// Version4 represents socks4
	Version4 Version = 0x04
	// Version5 represents socks5
	Version5 Version = 0x05
)

// RequestCmd is the requested socks command
type RequestCmd uint8

const (
	// RequestCmdConnect represents the CONNECT command
	RequestCmdConnect RequestCmd = 0x01
	// RequestCmdBind represents the BIND command
	RequestCmdBind RequestCmd = 0x02
	// RequestCmdAssociate represents the ASSOCIATE command
	RequestCmdAssociate RequestCmd = 0x03
)

// RequestAddressType is the Address Type from the socks communications
type RequestAddressType uint8

const (
	// RequestAddressTypeIPv4 represents IPv4
	RequestAddressTypeIPv4 RequestAddressType = 0x01
	// RequestAddressTypeDomainname represents a domain name
	RequestAddressTypeDomainname RequestAddressType = 0x03
	// RequestAddressTypeIPv6 represents IPv6
	RequestAddressTypeIPv6 RequestAddressType = 0x04
)

// Value gets the real value of the RequestAddressType
func (t RequestAddressType) Value() uint8 {
	return uint8(t)
}

// RequestReplyReason is used in replies to the client
type RequestReplyReason uint8

// Value gets the real value of the RequestReplyReason
func (r RequestReplyReason) Value() uint8 {
	return uint8(r)
}

const (
	// RequestReplySucceeded represents the "succeeded" reply
	RequestReplySucceeded RequestReplyReason = 0x00
	// RequestReplyGeneralFailure represents the "general SOCKS server failure" reply
	RequestReplyGeneralFailure RequestReplyReason = 0x01
	// RequestReplyConnectionNotAllowed represents the "connection not allowed by ruleset" reply
	RequestReplyConnectionNotAllowed RequestReplyReason = 0x02
	// RequestReplyNetworkUnreachable represents the "Network unreachable" reply
	RequestReplyNetworkUnreachable RequestReplyReason = 0x03
	// RequestReplyHostUnreachable represents the "Host unreachable" reply
	RequestReplyHostUnreachable RequestReplyReason = 0x04
	// RequestReplyConnectionRefused represents the "Connection refused" reply
	RequestReplyConnectionRefused RequestReplyReason = 0x05
	// RequestReplyTTLExpired represents the "TTL expired" reply
	RequestReplyTTLExpired RequestReplyReason = 0x06
	// RequestReplyCommandNotSupported represents the "Command not supported" reply
	RequestReplyCommandNotSupported RequestReplyReason = 0x07
	// RequestReplyAddressTypeNotSupported represents the "Address type not supported" reply
	RequestReplyAddressTypeNotSupported RequestReplyReason = 0x08
	// RequestReplyMethodNotSupported represents the "Method not supported" reply
	RequestReplyMethodNotSupported RequestReplyReason = 0xff
)

// RequestReply is the struct for the client reply
type RequestReply struct {
	Version     Version
	Reply       RequestReplyReason
	AddressType RequestAddressType
	BindAddress string
	BindPort    uint16
}

// Error is used to also return a ReplyReason to the client
type Error struct {
	Err    error
	Reason RequestReplyReason
}

// Error returns the underying error string
func (e *Error) Error() string { return e.Err.Error() }

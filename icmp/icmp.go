package icmp

import (
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// SendDstUnreach send a icmp (dest port unreachable) packet
func SendDstUnreach(laddr, raddr *net.IPAddr, udpPayload []byte) error {
	conn, err := net.DialIP("ip4:icmp", laddr, raddr)
	if err != nil {
		return err
	}

	body := icmp.DstUnreach{
		Data: udpPayload,
	}

	msg := icmp.Message{
		Type: ipv4.ICMPTypeDestinationUnreachable,
		Code: 3,
		Body: &body,
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return err
	}

	_, err = conn.Write(msgBytes)

	return err
}

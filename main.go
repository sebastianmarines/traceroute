package main

import (
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {

	packetconn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		panic(err)
	}

	defer packetconn.Close()

	message := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("Hello World"),
		},
	}

	messageBytes, err := message.Marshal(nil)
	if err != nil {
		panic(err)
	}

	ip := net.ParseIP("1.1.1.1")

	_, err = packetconn.WriteTo(messageBytes, &net.IPAddr{IP: ip})

	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1500)
	n, addr, err := packetconn.ReadFrom(buffer)

	if err != nil {
		panic(err)
	}

	message, err = icmp.ParseMessage(1, buffer[:n])

	if err != nil {
		panic(err)
	}

	switch message.Type {
	case ipv4.ICMPTypeEchoReply:
		echoReply := message.Body.(*icmp.Echo)
		println("Echo reply from", addr.String(), "with ID", echoReply.ID, "and Seq", echoReply.Seq)
	default:
		println("Unknown ICMP message type")
	}

}

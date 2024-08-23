package main

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {

	ip := net.ParseIP("1.1.1.1")

	packetconn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		panic(err)
	}

	defer packetconn.Close()

	// Send 64 messages
	for i := 0; i < 64; i++ {
		message := &icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   1,
				Seq:  i + 1,
				Data: []byte("Hello World"),
			},
		}
		messageBytes, err := message.Marshal(nil)
		if err != nil {
			panic(err)
		}

		if err := packetconn.IPv4PacketConn().SetTTL(i + 1); err != nil {
			panic(err)
		}

		if err := packetconn.IPv4PacketConn().SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
			panic(err)
		}

		start := time.Now()

		if _, err := packetconn.WriteTo(messageBytes, &net.IPAddr{IP: ip}); err != nil {
			panic(err)
		}

		buffer := make([]byte, 1500)
		n, addr, err := packetconn.ReadFrom(buffer)

		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			continue
		}

		elapsed := time.Since(start)

		message, err = icmp.ParseMessage(1, buffer[:n])

		if err != nil {
			panic(err)
		}

		switch message.Type {
		case ipv4.ICMPTypeEchoReply:
			echoReply := message.Body.(*icmp.Echo)
			fmt.Printf("Echo reply from %s with ID %d and Seq %d, time elapsed: %v\n", addr.String(), echoReply.ID, echoReply.Seq, elapsed)
			// Break the loop
			i = 64
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("Time exceeded from %s, time elapsed: %v\n", addr.String(), elapsed)
		case ipv4.ICMPTypeDestinationUnreachable:
			fmt.Printf("Destination unreachable from %s, time elapsed: %v\n", addr.String(), elapsed)
		default:
			fmt.Printf("Unknown message type: %v\n", message.Type)
		}

	}
}

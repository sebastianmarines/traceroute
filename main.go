package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <dest>")
		os.Exit(1)
	}

	host := ""

	// Parse the IP address
	ip := net.ParseIP(args[1])
	if ip == nil {
		// Try to resolve the IP address
		ip = getIp(args[1])
		if ip == nil {
			fmt.Println("Could not resolve the IP address")
			os.Exit(1)
		}
		host = os.Args[1]
	}

	if host == "" {
		host = getHost(ip)
	}

	maxHops := 64

	packetconn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		panic(err)
	}

	defer packetconn.Close()

	fmt.Printf("Traceroute to %s (%s), %d hops max\n", host, ip.String(), maxHops)

	for i := 0; i < maxHops; i++ {
		fmt.Printf("%2d ", i+1)

		lastAddress := ""

		for j := 0; j < 3; j++ {
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
				printErrorResponse(j)
				continue
			}

			elapsed := time.Since(start)

			message, err = icmp.ParseMessage(1, buffer[:n])

			if err != nil {
				panic(err)
			}

			host := getHost(addr.(*net.IPAddr).IP)

			switch message.Type {
			case ipv4.ICMPTypeEchoReply:
				printResponse(&lastAddress, host, addr.String(), elapsed, j)
				// Break the loop
				i = 64
			case ipv4.ICMPTypeTimeExceeded:
				printResponse(&lastAddress, host, addr.String(), elapsed, j)
			default:
				printErrorResponse(j)
			}

		}
	}
}

func getHost(ip net.IP) string {
	addrs, err := net.LookupAddr(ip.String())

	if err != nil {
		return ip.String()
	}

	if len(addrs) > 0 {
		return addrs[0]
	}

	return ip.String()
}

func printResponse(lastAddress *string, host string, addr string, elapsed time.Duration, j int) {
	if *lastAddress == "" {
		fmt.Printf("%s (%s) ", host, addr)
		*lastAddress = addr
	}

	if *lastAddress != addr {
		fmt.Printf("\n   %s (%s) ", host, addr)
		*lastAddress = addr
	}

	fmt.Printf("%s ", elapsed)

	if j == 2 {
		fmt.Println()
	}
}

func printErrorResponse(j int) {
	fmt.Printf("* ")
	if j == 2 {
		fmt.Println()
	}
}

func getIp(host string) net.IP {
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip
		}
	}

	return nil
}

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rnben/construct_icmp_packet/icmp"
	"github.com/rnben/construct_icmp_packet/udp"
)

func main() {
	go startUDPServer(30000)

	time.Sleep(time.Second * 2)

	// udp client
	conn, err := net.DialUDP("udp",
		&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 20000},
		&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 30000})
	if err != nil {
		fmt.Printf("dial server，err: %v\n", err)
		return
	}
	defer conn.Close()

	// keepalived heart
	go func() {
		for {
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				log.Printf("client send heart, err: %v\n", err)
			}

			time.Sleep(time.Second * 5)
		}
	}()

	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("Yeah server mock success, err: %v\n", err)
			continue
		}

		log.Printf("client recv: %s\n", string(data[:n]))
	}
}

func startUDPServer(port int) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		var data [1024]byte

		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("server recv: %s", data[:n])

		_, err = udpConn.WriteToUDP([]byte("pong"), addr)
		if err != nil {
			log.Printf("server send, err: %v\n", err)
		}

		err = icmp.SendDstUnreach(
			&net.IPAddr{IP: net.ParseIP("127.0.0.1")},
			&net.IPAddr{IP: net.ParseIP("127.0.0.1")},
			udp.ConstructUDPacket("127.0.0.1", "127.0.0.1", 20000, 30000),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

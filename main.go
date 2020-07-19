package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"time"
)

func main() {
	ssid := flag.String("s", "", "`SSID` to encode")
	password := flag.String("p", "", "`password` to encode")
	repeat := flag.Uint("repeat", 20, "Number of times to repeat the cofig sequence")
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	print := func(a interface{}) {
		if !*verbose {
			return
		}
		fmt.Print(a)
	}

	ips, err := encode(*ssid, *password)
	if err != nil {
		log.Fatal(err)
	}

	print(fmt.Sprintf("sending config sequence %d times\n", *repeat))
	for i := 0; i < int(*repeat); i++ {
		time.Sleep(100 * time.Millisecond)
		print(".")
		for j, ip := range ips {
			if err := sendTo(ip, []byte{'a'}); err != nil {
				log.Fatalf("failed to send packet %d. %s: %v", j+1, ip, err)
			}
		}
	}

	print(" done!\n")
}

func sendTo(ip net.IP, payload []byte) error {
	c, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: ip, Port: 1234})
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Write(payload)
	return err
}

func encode(ssid string, password string) ([]net.IP, error) {
	ls, lp := len(ssid), len(password)
	if ls == 0 || lp == 0 {
		return nil, fmt.Errorf("ssid and password are both required")
	}

	ips := make([]net.IP, 0, 12+((ls+lp)/2))

	ips = append(ips,
		// prefix
		net.ParseIP("226.0.1.2"),
		net.ParseIP("226.1.2.3"),
		net.ParseIP("226.2.3.4"),

		// len password
		net.ParseIP(fmt.Sprintf("226.32.%d.%d", lp, lp)),
	)

	// password bytes
	encoded := make([]byte, lp)
	for i := 0; i < lp; i++ {
		encoded[i] = password[i] ^ byte(0x50+i)
		i++
		b := byte(0)
		if i < lp {
			b = password[i] ^ byte(0x50+i)
			encoded[i] = b
		}
		ips = append(ips,
			net.ParseIP(fmt.Sprintf("226.%d.%d.%d", 64+(i/2), b, encoded[i-1])),
		)
	}

	// password checksum
	i := crc32.ChecksumIEEE(encoded)
	ips = append(ips,
		net.ParseIP(fmt.Sprintf("226.96.%d.%d", i&0x0000ff00>>8, i&0x000000ff)),
		net.ParseIP(fmt.Sprintf("226.97.%d.%d", i&0xff000000>>24, i&0x00ff0000>>16)),
	)

	// len ssid
	ips = append(ips,
		net.ParseIP(fmt.Sprintf("226.16.%d.%d", ls, ls)),
	)

	// ssid bytes
	for i := 0; i < ls; i++ {
		a := ssid[i]
		i++
		b := byte(0)
		if i < ls {
			b = ssid[i]
		}
		ips = append(ips,
			net.ParseIP(fmt.Sprintf("226.%d.%d.%d", 48+(i/2), int(b), int(a))),
		)
	}

	// ssid checksum
	i = crc32.ChecksumIEEE([]byte(ssid))
	ips = append(ips,
		net.ParseIP(fmt.Sprintf("226.80.%d.%d", i&0x0000ff00>>8, i&0x000000ff)),
		net.ParseIP(fmt.Sprintf("226.81.%d.%d", i&0xff000000>>24, i&0x00ff0000>>16)),
	)

	// FIN
	ips = append(ips,
		net.ParseIP("226.112.35.35"),
	)

	return ips, nil
}

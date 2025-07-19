package util

import (
	"fmt"
	"net"
)

func GetIP(real bool) string {
	if real {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			fmt.Println("Error dialing:", err)
			return "localhost"
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	} else {
		return "localhost"
	}
}

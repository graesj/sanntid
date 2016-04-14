package UPD

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	. "../.././message"
)

const (
	PORT = ":20001"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
		return true
	}
	return false
}

func UDPsend(channel chan Message) bool {

	Addr := []string{"129.241.187.255", PORT}
	broadcastUDP, err := net.ResolveUDPAddr("udp", strings.Join(Addr, ""))
	if CheckError(err) {
		return false
	}

	broadcastConn, err := net.DialUDP("udp", nil, broadcastUDP)
	if CheckError(err) {
		return false
	}

	defer broadcastConn.Close()

	for {
		buf, err := json.Marshal(<-channel)

		if !CheckError(err) {
			broadcastConn.Write(buf)
		}
	}
}

func UDPlisten(channel chan Message) bool {

	//Connect to network
	UDPRecAddr, err := net.ResolveUDPAddr("udp", PORT)
	if CheckError(err) {
		return false
	}

	UDPConn, err := net.ListenUDP("udp", UDPRecAddr)
	if CheckError(err) {
		return false
	}

	//Close connection after receiving is complete
	defer UDPConn.Close()

	//HEY! LISTEN NA'VI
	buf := make([]byte, 2048)
	trimmed_buf := make([]byte, 1)
	var received_message Message

	for {
		n, _, err := UDPConn.ReadFromUDP(buf)
		CheckError(err)
		trimmed_buf = buf[:n]
		err = json.Unmarshal(trimmed_buf, &received_message)
		if err == nil {
			channel <- received_message
		}

	}
	return true
}

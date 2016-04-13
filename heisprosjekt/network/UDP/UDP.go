package UPD

import (
	. "../.././message"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	PORT = ":20001"
)

//Sending and receiving data from UDP-multicast-network
func CheckError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
		return true

	}
	return false
}

func UDPsend(channel chan Message) bool {

	//Connect to UDP network
	Addr := []string{"129.241.187.255", PORT}
	broadcastUDP, err := net.ResolveUDPAddr("udp", strings.Join(Addr, ""))
	if CheckError(err) {
		return false
	}

	broadcastConn, err := net.DialUDP("udp", nil, broadcastUDP)
	if CheckError(err) {
		return false
	}

	//Close connection after packet is sent/or fails to send
	defer broadcastConn.Close()
	//Send packet

	for {
		buf, err := json.Marshal(<-channel) //JavaScript Object Notation, used for data
		//interchanging
		//Channels for semaphore
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
	fmt.Println("Shi")
	return true
}

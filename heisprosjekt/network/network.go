package network

import (
	. "./UDP"
	"time"
	//"net"
	. ".././message"
	//"fmt"
)

func ip_broadcast(ip_key int, UDPsend chan Message) {

	UDPsend <- Message{IP: ip_key, ID: 1}
	time.Sleep(100 * time.Millisecond)

}

func Manager(fromMain chan Message, toMain chan Message) {

	ip_key := int(addr[1].String()[12] - '0') * 100 + int(addr[1].String()[13] - '0') * 10 + int(addr[1].String()[14] - '0')

	sendChan := make(chan Message, 50)
	recieveChan := make(chan Message, 50)

	go ip_broadcast(ip_key, sendChan)
	go UDPsend(sendChan)
	go UDPlisten(recieveChan)

	for {
		select {
		case message := <-recieveChan:

			if message.ID = IP
				

		case message := <-fromMain:

			sendChan <- message
		}
	}

}

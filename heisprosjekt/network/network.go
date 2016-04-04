package network

import (
	. "./UDP"
	"time"
	//"net"
	. ".././message"
	//"fmt"
)

func ip_broadcast(ip int, UDPsend chan Message) {

	UDPsend <- Message{IP: ip, ID: 1}
	time.Sleep(100 * time.Millisecond)

}

func Manager(fromMain chan Message, toMain chan Message) {

	IP := 10

	sendChan := make(chan Message, 50)
	recieveChan := make(chan Message, 50)

	go ip_broadcast(IP, sendChan)
	go UDPsend(sendChan)
	go UDPlisten(recieveChan)

	for {
		select {
		case message := <-recieveChan:

			toMain <- message

		case message := <-fromMain:

			sendChan <- message
		}
	}

}

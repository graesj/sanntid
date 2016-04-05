package network

import (
	. "./UDP"
	"time"
	//"net"
	. ".././message"
	//"fmt"
)

func ip_broadcast(ip_key int, UDPsend chan Message) {

	UDPsend <- Message{Source: ip_key, Id: ID}
	time.Sleep(100 * time.Millisecond)

}

func Manager(fromMain chan Message, toMain chan Message) {
	addr, _ := net.InterfaceAddrs()
	ip_key := int(addr[1].String()[12] - '0') * 100 + int(addr[1].String()[13] - '0') * 10 + int(addr[1].String()[14] - '0')

	sendChan := make(chan Message, 50)
	recieveChan := make(chan Message, 50)

	go ip_broadcast(ip_key, sendChan)
	go UDPsend(sendChan)
	go UDPlisten(recieveChan)


	con_timer := make(map[int]*time.Timer)



	for {
		select {
		case message := <-recieveChan:

			if message.ID = IP {
				if (con_timer[message.IP] != 0){ //The ip_key already has a running Timer
					con_timer[message.IP].Reset(3*time.Second)

				}
				else { //new elevator
					con_timer[message.IP].AfterFunc(3*time.Second, remove_elev(message.IP))
					message.ID := NEW_ELEVATOR
					
					toMain <- message
				}
			}

			else{
				toMain <- message
			}

				

		case message := <-fromMain:

			sendChan <- message
		}
	}

}


func remove_elev(ip_key int){
	m := Message{IP: ip_key, ID = REMOVE_ELEVATOR}
	toMain <- m
	delete(con_timer,ip_key)
}
package network

import (
	. "../message"
	. "./UDP"
	"time"
	"net"
	//"fmt"
//	. ".././elev_manager"
	//. ".././elev_manager/fsm"
	. "../structs"
)

var con_timer map[int]*time.Timer

func BroadcastElevatorInfo(e Elevator, UDPsend chan Message) {

	UDPsend <- Message{Source: e.Self_id, ID: ELEVATOR_DATA, Elevator: e}
}

func Manager(fromMain chan Message, toMain chan Message) {

	sendChan := make(chan Message, 50)
	recieveChan := make(chan Message, 50)

	go UDPsend(sendChan)
	go UDPlisten(recieveChan)

	con_timer = make(map[int]*time.Timer)

	for {
		select {
		case message := <-recieveChan:

			if message.ID == ELEVATOR_DATA {
				_, present := con_timer[message.Source]

				if present { //The ip_key already has a running Timer
					con_timer[message.Source].Reset(500*time.Millisecond)
					toMain <- message
				} else { //new elevator
					con_timer[message.Source] = time.AfterFunc(500*time.Millisecond, func() { remove_elev(message.Source, toMain) })
					message.ID = NEW_ELEVATOR

					toMain <- message
				}

			} else {
				toMain <- message
			}

		case message := <-fromMain:

			sendChan <- message
		}
	}

}

func remove_elev(ip_key int, toMain chan Message) {
	m := Message{Source: ip_key, ID: REMOVE_ELEVATOR}
	toMain <- m
	delete(con_timer, ip_key)
}
func GetLastNumbersOfIp() int {
	addr, _ := net.InterfaceAddrs()
	ip := int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')

	return ip
}
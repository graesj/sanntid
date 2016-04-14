package network

import (
	"net"
	"time"

	. "../message"
	. "../structs"
	. "./UDP"
)

var con_timer map[int]*time.Timer

func BroadcastElevatorInfo(e Elevator, fromMain chan Message) {

	fromMain <- Message{Source: e.Self_id, ID: ELEVATOR_UPDATE, Elevator: e}
}

func NetworkManager(fromMain chan Message, toMain chan Message) {

	sendChan := make(chan Message, 50)
	receiveChan := make(chan Message, 50)

	go UDPsend(sendChan)
	go UDPlisten(receiveChan)

	con_timer = make(map[int]*time.Timer)

	for {
		select {
		case message := <-receiveChan:

			if message.ID == ELEVATOR_UPDATE {
				_, present := con_timer[message.Source]

				if present {
					con_timer[message.Source].Reset(500 * time.Millisecond)
					toMain <- message
				} else {
					con_timer[message.Source] = time.AfterFunc(500*time.Millisecond, func() { remove_elev(message.Source, toMain, message.Elevator) })
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

func remove_elev(ip_key int, toMain chan Message, elev Elevator) {
	elev.ErrorType = ERROR_NETWORK
	m := Message{Source: ip_key, ID: REMOVE_ELEVATOR, Elevator: elev}
	toMain <- m
	delete(con_timer, ip_key)
}
func GetLastNumbersOfIp() int {
	addr, _ := net.InterfaceAddrs()
	ip := int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')

	return ip
}

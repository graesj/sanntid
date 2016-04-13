package main

import (
	. "./elev_manager"
	. "./elev_manager/fsm/driver"
	. "./message"
	. "./network"
	. "./structs"
	//. "./utilities"
	//. "./elev_manager/fsm"
	. "fmt"
	"time"
)

func main() {

	e := MakeElevManager()
	buttonChan := make(chan Message, 100)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go e.ProcessElevOrders(fromMain)
	go Manager(fromMain, toMain)
	go CheckButtons(buttonChan)

	broadcastTicker := time.NewTicker(100 * time.Millisecond).C

	for {

		select {
		case message := <-toMain:

			switch message.ID {

			case REMOVE_ELEVATOR:

				e.Elevators[message.Source].ErrorType = message.Elevator.ErrorType
				Print("Feil oppdaget av type: ")
				Println(e.Elevators[message.Source].ErrorType)
				e.OnConnectionTimeout(message.Source, fromMain, message.Elevator)

			case BUTTON_EXTERNAL:
				ElevSetButtonLamp(message.ButtonType, message.Floor, 1)
				if e.IsMaster() {
					Println("Mottok knapp og er master....")
					assignID := e.DetermineTargetElev(message.ButtonType, message.Floor)
					message.Source = e.Self_id
					message.ID = ORDER_COMMAND
					message.Target = assignID
					Print("Beste id ble beregnet til å være: ")
					Println(message.Target)
					fromMain <- message
					time.AfterFunc(1000*time.Millisecond, func() { e.CheckIfOrderIsReceived(message, fromMain) })

				}

			case NEW_ELEVATOR:
				_, present := e.Elevators[message.Source]

				if present && (e.Self_id != message.Source) {

					e.ResendInitialOrders(fromMain, message.Source)

				} else {
					e.NewElevator(message.Elevator)
				}
				e.Elevators[message.Source].ErrorType = ERROR_NONE
				e.UpdateMaster(-1)

			case ELEVATOR_DATA:
				if message.Elevator.Self_id == e.Self_id {
					//Println("Mottok egen melding")
				} else {
					//Println("Mottok annen heismelding...")
					_, present := e.Elevators[message.Elevator.Self_id]

					if present { //Update the elevatordata
						e.ElevatorUpdate(message.Elevator)
						//	Println("Det var en oppdatering")
					}
				}
			case ORDER_COMMAND:
				if message.Target == e.Self_id {
					Println("En ektern kommando fra master ble sent til meg:DDD")
					e.AddExternalOrders(message.Floor, message.ButtonType)
				}
			case LAMP_MESSAGE:
				ElevSetButtonLamp(message.ButtonType, message.Floor, 0)

			case GET_UP_TO_DATE:
				//This happens when an elevator has been disconnected. Gets resend of own internal orders
				if message.Target == e.Self_id {
					if e.Elevators[e.Self_id].ErrorType == ERROR_NETWORK {

					} else {
						e.CopyInternalOrder(message.Elevator)
					}
				}
				e.Elevators[e.Self_id].ErrorType = ERROR_NONE
				e.UpdateMaster(-1)
			}

		case buttonMessage := <-buttonChan:

			if buttonMessage.ID == BUTTON_EXTERNAL {
				fromMain <- buttonMessage

				//en bool som sier at ordre har blitt sent
				//kanskje ha en funksjon som sjekker at master plukker den opp
				Println("Sender eksterntknappetrykk ut")

			} else if buttonMessage.ID == BUTTON_INTERNAL {
				e.AddInternalOrders(buttonMessage.Floor, BTN_CMD)
				Println("bais")
			}

		case <-broadcastTicker:
			BroadcastElevatorInfo(*e.Elevators[e.Self_id], fromMain)
			if e.IsMaster() {
				//Println("Jeg er master!!!!")
			}

		}
	}
}

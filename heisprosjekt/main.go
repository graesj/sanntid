package main

import (
	. "fmt"
	"time"

	. "./elev_manager"
	. "./elev_manager/fsm"
	. "./elev_manager/fsm/driver"
	. "./message"
	. "./network"
	. "./structs"
)

func main() {

	e := MakeElevManager()
	buttonChan := make(chan Message, 100)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go NetworkManager(fromMain, toMain)
	go CheckButtons(buttonChan)
	time.AfterFunc(200*time.Millisecond, func() { go ProcessElevOrders(e.Elevators[e.Self_id], fromMain) })

	broadcastTicker := time.NewTicker(100 * time.Millisecond).C

	for {

		select {
		case message := <-toMain:

			switch message.ID {

			case REMOVE_ELEVATOR:

				if message.Source == e.Self_id {
					TurnOffAllExternalLights()

				}
				e.Elevators[message.Source].ErrorType = message.Elevator.ErrorType
				Print("Elevator has dropped out. ID: ")
				Println(message.Source)
				e.OnConnectionTimeout(message.Source, fromMain, message.Elevator)

			case BUTTON_EXTERNAL:

				ElevSetButtonLamp(message.ButtonType, message.Floor, 1)
				if e.IsMaster() {
					Println("Mottok knapp og er master. Beregner beste ID for oppdraget..")
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

					e.ResendInternalOrders(fromMain, message.Source)

				} else if !present {
					e.NewElevator(message.Elevator)
				}
				e.Elevators[message.Source].ErrorType = ERROR_NONE
				e.UpdateMaster(-1)

			case ELEVATOR_UPDATE:

				if message.Elevator.Self_id == e.Self_id {
					//Your own message. Do nothing
				} else {

					_, present := e.Elevators[message.Elevator.Self_id]

					if present {
						e.OnElevatorUpdate(message.Elevator)
					}
				}

			case ORDER_COMMAND:

				if message.Target == e.Self_id {
					AddExternalOrders(e.Elevators[e.Self_id], message.Floor, message.ButtonType)
				}

			case LAMP_MESSAGE:

				ElevSetButtonLamp(message.ButtonType, message.Floor, 0)

			case GET_UP_TO_DATE:
				//This happens when an elevator has been disconnected. Gets resend of own internal orders
				if message.Target == e.Self_id {
					if e.Elevators[e.Self_id].ErrorType == ERROR_NETWORK {

					} else {
						e.CopyInternalOrders(message.Elevator)
					}
				}
				e.Elevators[e.Self_id].ErrorType = ERROR_NONE
				e.UpdateMaster(-1)
			}

		case buttonMessage := <-buttonChan:

			if buttonMessage.ID == BUTTON_EXTERNAL {
				fromMain <- buttonMessage
				Println("Sender eksterntknappetrykk ut")

			} else if buttonMessage.ID == BUTTON_INTERNAL {
				AddInternalOrders(e.Elevators[e.Self_id], buttonMessage.Floor, BTN_CMD)
				Println("bais")
			}

		case <-broadcastTicker:
			BroadcastElevatorInfo(*e.Elevators[e.Self_id], fromMain)
			//Println(e.Elevators[Self_id].Orders)

		}
	}
}

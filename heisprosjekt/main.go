package main

import (
	. "./message"
	"./network"
	"./elev_manager"
	//"fmt"
	//"time"
)

func main() {
	e := Em_makeElevManager()
	buttonChan := make(chan Message)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go network.Manager(fromMain, toMain)
	go CheckButtons(fromMain, e)

	mes := Message{Source: 1, Floor: 1, Target: 1, ID: 1, IP: 1}
	i := 0

	for {
		i = i + 1
		select {
		case message := <-toMain:

			switch message.Id {

			case NEW_ELEVATOR:
				e.newElevator(message) //Skal legge til den nye heisen, og sjekke hvem som er master

			case REMOVE_ELEVATOR:
				e.removeElevator(message.Source)

			/*case GENERAL_UPDATE:
				if message.Source != e.Id {
					e.updateElevators(message)
				}
*/
			//case CALCULATE_COST:
			//	e.calculateCostOfOrder()

			//case BUTTON_EXTERNAL:
			//	if(e.isMaster()){
					//En funksjon som ber alle kalkulere kosten for å ta oppdraget.
			//	}
			}

		case mes <- toMain:
			fromMain <- mes
			mes.ID = i

		case buttonMessage := <- buttonChan:
			if buttonMessage.ID = BUTTON_EXTERNAL {
				fromMain <- button

			} else if buttonMessage.ID = BUTTON_INTERNAL {
				e.Em_handleFloorButtonPressed(buttonMessage.ButtonType,buttonMessage.Floor)
			}

		}
	}
}

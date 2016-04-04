package elev_manager

import (
	. "./fsm"
)

type elev_manager struct{
	master int
	self_id int
	external_orders[N_FLOORS*2-2]  //This is where the orders from the floor panels are put. These orders are broadcasted to all elevators, 
	elevators map[int] *Elevator //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func Em_checkButtons() int {
	return Hello()
}


/*
Check all buttons to see if any action must be performed, and initialize an action
*/
func Fsm_checkAllButtons(){
	for floor := 0; floor < 4; floor++ {
		for buttonType := 0; buttonType < 3; buttonType++{
			if ElevGetButtonSignal(buttonType, floor) {
				Fsm_handleButtonPressed(buttonType, floor)
			}
		}
	}
}


/*
Perform an action based on a button being pressed
*/
func (e * Elevator) Fsm_handleButtonPressed(e Elevator, buttonType int, floor int) { 


	if buttonType != 2 {
		//broadcast button order to the master. Also remember to send direction 'dir'.
	}

	else{
		switch floor {
		case 0:
			e.elevList[floor] = 1
		case 3:
			e.elevList[floor+2] = 1
		default:
			e.elevList[floor] = 1
			e.elevList[floor+2] = 1
		}
	}
}

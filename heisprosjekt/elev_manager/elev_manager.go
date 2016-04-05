package elev_manager

import (
	//. ".././message"
	//. ".././network"
	. "./fsm"
	. "./fsm/driver"
)

type elev_manager struct {
	master          int
	Self_id         int
	external_orders [N_FLOORS*2 - 2]int //This is where the orders from the floor panels are put. These orders are broadcasted to all Elevators,
	Elevators       map[int]*Elevator   //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func Em_makeElevManager() elev_manager {
	var e elev_manager
	e.Elevators = make(map[int]*Elevator)
	e.Self_id = 0
	e.Elevators[e.Self_id] = new(Elevator)
	e.Elevators[e.Self_id].Fsm_initiateElev()

	return e
}

/*
Check all buttons to see if any action must be performed, and initialize an action
*/
func (e *elev_manager) Em_checkAllFloorButtons() {
	for floor := 0; floor < 4; floor++ {
		for buttonType := 0; buttonType < 3; buttonType++ {
			if (floor == 0 && buttonType == 1) || (floor == N_FLOORS-1 && buttonType == 0) {

			} else {
				if ElevGetButtonSignal(buttonType, floor) == 1 {
					e.Em_handleFloorButtonPressed(buttonType, floor)
				}
			}
		}
	}
}

/*
Perform an action based on a button being pressed
*/
func (e *elev_manager) Em_handleFloorButtonPressed(buttonType int, floor int) {

	if buttonType != 2 {
		//broadcast button order to the master. Also remember to send direction 'dir'.
	} else {

		//The elevator panel has been used (inside the elevator), and the internal orders of this elevator should be updated.

		e.Elevators[e.Self_id].Internal_orders[0][floor] = 1
		e.Elevators[e.Self_id].Internal_orders[1][floor] = 1

	}
}

func (e *elev_manager) Em_processElevOrders() {
	check4DirChange := 1
	switch e.Elevators[e.Self_id].State {
	case STATE_IDLE:
		e.Elevators[e.Self_id].State = STATE_RUNNING
		//Send message to master to inform him that you are available
	case STATE_RUNNING:
		switch e.Elevators[e.Self_id].Dir {
		case DIR_UP:
			for floor := e.Elevators[e.Self_id].Floor; floor < N_FLOORS; floor++ {
				if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 {
					check4DirChange = 0
					ElevSetMotorDirection(DIR_UP)
					if ElevGetFloorSensorSignal() == floor {
						StopAndOpenDoor()
						e.Elevators[e.Self_id].Floor = floor
						e.Em_UpdateInternalOrders(floor)
					}
				}
			}
			if check4DirChange == 1 {
				e.Elevators[e.Self_id].Dir = DIR_STOP
			}
		case DIR_DOWN:
			for floor := e.Elevators[e.Self_id].Floor; floor >= 0; floor-- {
				if e.Elevators[e.Self_id].Internal_orders[1][floor] == 1 {
					check4DirChange = 0
					ElevSetMotorDirection(DIR_DOWN)
					if ElevGetFloorSensorSignal() == floor {
						StopAndOpenDoor()
						e.Elevators[e.Self_id].Floor = floor
						e.Em_UpdateInternalOrders(floor)
					}
				}
			}
			if check4DirChange == 1 {
				e.Elevators[e.Self_id].Dir = DIR_STOP
			}
		case DIR_STOP:
			check4DirChange = 1
			for floor := e.Elevators[e.Self_id].Floor; floor < N_FLOORS; floor++ {
				if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 {
					e.Elevators[e.Self_id].Dir = DIR_UP
					break
				}
			}
			for floor := e.Elevators[e.Self_id].Floor; floor >= 0; floor-- {
				if e.Elevators[e.Self_id].Internal_orders[1][floor] == 1 {
					e.Elevators[e.Self_id].Dir = DIR_DOWN
					break
				}
			}
		}
	}
}

func (e *elev_manager) Em_UpdateInternalOrders(floor int) {
	e.Elevators[e.Self_id].Internal_orders[0][floor] = 0
	e.Elevators[e.Self_id].Internal_orders[1][floor] = 0
}

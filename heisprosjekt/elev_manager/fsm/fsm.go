package fsm

import (
	. "./driver"
	"fmt"
)

const (
	//STATES (Golang doesn't have enums)
	STATE_IDLE     = 0
	STATE_RUNNING  = 1
	STATE_DOOROPEN = 2
	STATE_STOP     = 3
)

type Elevator struct {
	State           int
	Dir             int
	Floor           int
	Internal_orders [2][N_FLOORS]byte //both external and internal orders
	External_orders [2][N_FLOORS]byte //orders from the external panel
	//Just for backup
}

/*		 floor:  0  1  2  3		dir:
  internal_orders[6] =   [*  *  *		up
			   *  *  *];	down
*/

/*
Creates an elevator struct object. Initializes the elevator: Runs to first floor, and sets all elevator parameters.
*/

func (e Elevator) Fsm_initiateElev() {
	ElevInit() //This function is necessary to reset all hardware

	if ElevGetFloorSensorSignal() == -1 {
		e.Dir = DIR_UP
		ElevSetMotorDirection(DIR_UP)
		for ElevGetFloorSensorSignal() == -1 {
		}
	}

	ElevSetMotorDirection(DIR_STOP)

	if ElevGetFloorSensorSignal() != 0 {
		e.Dir = DIR_DOWN
		ElevSetMotorDirection(DIR_DOWN)
		for ElevGetFloorSensorSignal() != 0 {
		}
	}
	ElevSetMotorDirection(DIR_STOP)
	e.Dir = DIR_STOP
	e.Floor = 0
	e.State = STATE_IDLE

	for i := 0; i < N_FLOORS; i++ {
		e.Internal_orders[0][i] = 0
		e.Internal_orders[1][i] = 0
	}
	fmt.Println("Initialized")
}

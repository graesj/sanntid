package fsm

import (
	. "./driver"
	"fmt"
)

const (
	N_FLOORS = 4
	UP       = 1
	DOWN     = -1
	STOP     = 0

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
	Internal_orders [N_FLOORS*2 - 2]byte
}

/*		 floor:  0  1  2  3		dir:
  internal_orders[6] =   [*  *  *		up
			   *  *  *];	down
*/

/*
Creates an elevator struct object. Initializes the elevator: Runs to first floor, and sets all elevator parameters.
*/
func Hello() int {
	fmt.Println("Hello")
	return 0
}

func (e Elevator) Fsm_initiateElev() {
	ElevInit() //This function is necessary to reset all hardware

	if ElevGetFloorSensorSignal() == -1 {
		e.Dir = UP
		ElevSetMotorDirection(UP)
		for ElevGetFloorSensorSignal() == -1 {
		}
	}

	ElevSetMotorDirection(STOP)

	if ElevGetFloorSensorSignal() != 0 {
		e.Dir = DOWN
		ElevSetMotorDirection(DOWN)
		for ElevGetFloorSensorSignal() != 0 {
		}
	}
	ElevSetMotorDirection(STOP)
	e.Dir = STOP
	e.Floor = 0
	e.State = STATE_IDLE

	for i := 0; i < N_FLOORS*2-2; i++ {
		e.Internal_orders[i] = 0
	}
	fmt.Println("Initialized")
}

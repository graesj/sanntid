package fsm

import (
	. "./driver"
	"fmt"
	. "../.././structs"
	//. "../"
)

const (
	//STATES (Golang doesn't have enums)
	STATE_IDLE     = 0
	STATE_RUNNING  = 1
	STATE_DOOROPEN = 2
	STATE_STOP     = 3
)





/*		 floor:  0  1  2  3		dir:
  internal_orders[6] =   [*  *  *		up
			   *  *  *];	down
*/

/*
Creates an elevator struct object. Initializes the elevator: Runs to first floor, and sets all elevator parameters.
*/

func (e *Elevator) Fsm_initiateElev() {
	ElevInit() //This function is necessary to reset all hardware

	if ElevGetFloorSensorSignal() == -1 {
		//e.Dir = DIR_UP
		//ElevSetMotorDirection(DIR_UP)
		for ElevGetFloorSensorSignal() == -1 {
			fmt.Println(ElevGetFloorSensorSignal())
		}
	}
	fmt.Println("KJOR OPP")
	ElevSetMotorDirection(DIR_STOP)

	if ElevGetFloorSensorSignal() != 0 {
		//e.Dir = DIR_DOWN
		
		for ElevGetFloorSensorSignal() != 0 {
			ElevSetMotorDirection(DIR_DOWN)
		}
	}
	ElevSetMotorDirection(DIR_STOP)
	/*
	e.Dir = DIR_STOP
	e.Floor = 0
	e.State = STATE_IDLE

	for i := 0; i < N_FLOORS; i++ {
		e.Internal_orders[0][i] = 0
		e.Internal_orders[1][i] = 0
	}
	*/
	fmt.Println("Initialized")
}

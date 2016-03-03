package fsm

import (
	. "./driver"
)

//#define up 0
//#define down 1

const (
	N_FLOORS = 4;
	UP = 0;
	DOWN = 1;

	//STATES (Golang doesn't have enums)
	FSM_IDLE = 0;
	FSM_RUNNING = 1;
	FSM_DOOROPEN = 2;
	FSM_STOP = 3;	

)



type Elevator struct{
	State int
	Dir int
	Floor int
	elevList[N_FLOORS*2-2] byte
}

/*		 floor:  0  1  2  3		dir:
	  elevList[6] = [*  *  *		up
						*  *  *];	down
	*/

func fsm_makeElev(){  
	var e Elevator
	
	//run to first floor

	e.Floor = 0
	e.Dir = 0
	e.State = FSM_IDLE
	for i := 0; i < N_FLOORS*2-2; i++ {
		e.elevList[i] = 0
	}
}

func fsm_buttonPressed(e Elevator, buttonType int, floor int){
	
	switch e.State {
		case FSM_IDLE:
			fsm_checkButton(e, buttonType, floor)
			//put on display lights
			e.State = FSM_RUNNING
			break
		case FSM_RUNNING:
			fsm_checkButton(e, buttonType, floor)
			//put on display lights
	    	break
		case FSM_DOOROPEN:
			if (floor != e.Floor){
				fsm_checkButton(e, buttonType, floor)
				//put on display lights
			}
			break
		case FSM_STOP:
			//alert the janitor, and start playing music for the stuck passengers...
			break 
	}
}

func Test() int{
	return ElevGetButtonSignal(0, 0)
}



func fsm_checkButton(e Elevator, buttonType int, floor int){

	if (buttonType != 2){
		//broadcast button order to master elevator. Also remember to send direction 'dir'.
	}

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




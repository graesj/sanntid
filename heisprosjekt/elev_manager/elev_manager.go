package elev_manager

import (
	. ".././message"
	//. ".././network"
	. "./fsm"
	. "./fsm/driver"
	. "fmt"
	"time"
)

type elev_manager struct {
	master          bool 
	Self_id         int
	external_orders [N_FLOORS*2 - 2]int //This is where the orders from the floor panels are put. These orders are broadcasted to all Elevators,
	Elevators       map[int]*Elevator   //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func Em_makeElevManager() elev_manager {
	var e elev_manager
 	addr, _ := net.InterfaceAddrs()
	e.Self_id = int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')
	e.Elevators = make(map[int]*Elevator)
	e.Self_id = 0
	e.Elevators[e.Self_id] = new(Elevator)
	e.Elevators[e.Self_id].Fsm_initiateElev()
	e.master = true

	return e
}

func (e *elev_manager) Em_newElevator(message Message){

	e.Elevators[message.SELF_ID] = Elevator{}

	if message.SELF_ID < e.Self_id {
		e.master = false
	}

}

func (e *elev_manager) Em_processElevOrders() {
	for {
		Println(e.Elevators[e.Self_id].Internal_orders)
		//Println(e.Elevators[e.Self_id].State)
		check4DirChange := 1
		switch e.Elevators[e.Self_id].State {
		case STATE_IDLE:
			e.Elevators[e.Self_id].State = STATE_RUNNING
			//Send message to master to inform him that you are available
			/*for  floor := 0, floor < 3, floor++ {
				if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 || e.Elevators[e.Self_id].Internal_orders[1][floor] == 1 {
					e.Elevators[e.Self_id].State = STATE_RUNNING
				}
			}*/
		case STATE_RUNNING:
			switch e.Elevators[e.Self_id].Dir {
			case DIR_UP:
				for floor := e.Elevators[e.Self_id].Floor; floor < N_FLOORS; floor++ {
					if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 {
						check4DirChange = 0
						ElevSetMotorDirection(DIR_UP)
						if ElevGetFloorSensorSignal() == floor {
							e.StopAndOpenDoor(floor)
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
							e.StopAndOpenDoor(floor)
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
		case STATE_DOOROPEN:
			ElevSetMotorDirection(DIR_STOP)

		}
	}
}

func (e *elev_manager) doorTimeout(floor int) {
	ElevSetDoorOpenLamp(0)
	e.Em_RemoveOrders(floor)

	e.Elevators[e.Self_id].State = STATE_RUNNING
}

func (e *elev_manager) Em_RemoveOrders(floor int) { //HEI BEDRE NAVN DA

	if e.Elevators[Self_if].Dir == DIR_DOWN
	e.Elevators[e.Self_id].Internal_orders[0][floor] = 0
	e.Elevators[e.Self_id].Internal_orders[1][floor] = 0

	if e.Elevators[e.Self_id].Dir == DIR_DOWN {
		e.Elevators[e.Self_id].External_orders[1][floor] = 1
	} else if e.Elevators[e.Self_id].Dir == DIR_UP {
		e.Elevators[e.Self_id].External_orders[0][floor] = 1
	}

}

func (e *elev_manager) Em_AddInternalOrders(floor int) {
	e.Elevators[e.Self_id].Internal_orders[0][floor] = 1
	e.Elevators[e.Self_id].Internal_orders[1][floor] = 1
}

func (e *elev_manager) Em_AddExternalOrders(floor int, buttonType int) {
	if buttonType == BTN_UP {
		if floor == 0 {
			e.Elevators[e.Self_id].Internal_orders[0][floor] = 1
			e.Elevators[e.Self_id].Internal_orders[1][floor] = 1

			e.Elevators[e.Self_id].External_orders[0][floor] = 1
			e.Elevators[e.Self_id].External_orders[1][floor] = 1
		} else {
			e.Elevators[e.Self_id].Internal_orders[0][floor] = 1
			e.Elevators[e.Self_id].External_orders[0][floor] = 1
		}
	} else if buttonType == BTN_DOWN {
		if floor == 3 {
			e.Elevators[e.Self_id].Internal_orders[0][floor] = 1
			e.Elevators[e.Self_id].Internal_orders[1][floor] = 1

			e.Elevators[e.Self_id].External_orders[0][floor] = 1
			e.Elevators[e.Self_id].External_orders[1][floor] = 1
		} else {
			e.Elevators[e.Self_id].Internal_orders[1][floor] = 1
			e.Elevators[e.Self_id].External_orders[1][floor] = 1
		}
	}
}

func (e *elev_manager) Em_NewElevator(elevMessage Message) {

	//e.Elevators[elevMessage.Source] = Elevator{}

}

func (e *elev_manager) StopAndOpenDoor(floor int) { //needs a better name
	ElevSetMotorDirection(DIR_STOP)
	Println("stopping")
	e.Elevators[e.Self_id].Floor = floor
	ElevSetDoorOpenLamp(1)
	time.AfterFunc(time.Second*3, func() { e.doorTimeout(floor) })
	e.Elevators[e.Self_id].State = STATE_DOOROPEN

}

func (e *elev_manager) NewElevator(msg Message) {

}

func (e *elev_manager) RemoveElevator(target int) {
	if e.Self_id == target {

	}
}

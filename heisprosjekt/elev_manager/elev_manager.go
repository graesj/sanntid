package elev_manager

import (
	. ".././message"
	. ".././network"
	. "./fsm"
	. "./fsm/driver"
	. "fmt"
	. ".././structs"
	"time"
	"math"
)

type elev_manager struct {
	master          int 
	Self_id         int
	external_orders [2][N_FLOORS]int //This is where the orders from the floor panels are put. These orders are broadcasted to all Elevators,
	Elevators       map[int]*Elevator   //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func Em_makeElevManager() elev_manager {
	var e elev_manager
	e.Elevators = make(map[int]*Elevator)
	e.Self_id = GetLastNumbersOfIp()
	e.Elevators[e.Self_id] = new(Elevator)
	Print("FØR")
	//Fsm_initiateElev()
	e.Elevators[e.Self_id].Fsm_initiateElev()
	Println("ETTER")
	e.Elevators[e.Self_id].Current_Dir = DIR_STOP
	e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
	e.Elevators[e.Self_id].Floor = 0
	e.Elevators[e.Self_id].State = STATE_IDLE
	e.Elevators[e.Self_id].Self_id = e.Self_id
	e.master = e.Self_id

	return e
}

func (e *elev_manager) Em_newElevator(elev Elevator){


	e.Elevators[elev.Self_id] = &elev
	Print("Ny heis er lagt til med self_id:")
	Println(elev.Self_id)
	Println("Lagt inn i lista: ", e.Elevators[elev.Self_id].Self_id)
	if elev.Self_id < e.Self_id {
		e.master = elev.Self_id
		Print("Den nye heisen ble master")
	}
}

func (e *elev_manager) Em_handleExternalOrder(buttonType int, floor int) int{
	bestID := -1
	lowestCost := 1000
	for key, elev:= range e.Elevators {
		curCost := calculateCost(*elev, buttonType, floor)
		if curCost < lowestCost {
			lowestCost = curCost
			bestID = key
		}
	}

	return bestID
}

func (e *elev_manager) Em_isMaster() bool {
	if (e.master == e.Self_id){
		return true
	} else {
		return false
	}
}

func (e *elev_manager) Em_elevatorUpdate(elev Elevator) {
	e.Elevators[elev.Self_id] = &elev
}

func (e *elev_manager) Em_processElevOrders() {
	Current_floor := -1
	furthest_floor := -1

	for {
		//Println(e.Elevators[e.Self_id].Internal_orders)
		//Println(e.Elevators[e.Self_id].State)
		switch e.Elevators[e.Self_id].State {
		case STATE_IDLE:
			e.Elevators[e.Self_id].State = STATE_RUNNING
			//Send message to master to inform him that you are available
			for floor := 0; floor < 3; floor++ {
				if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 || e.Elevators[e.Self_id].Internal_orders[1][floor] == 1 {
					e.Elevators[e.Self_id].State = STATE_RUNNING
				}
			}
		case STATE_RUNNING:
				switch e.Elevators[e.Self_id].Current_Dir {
				case DIR_UP:
					ElevSetMotorDirection(DIR_UP)
					if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
						for furthest_floor = N_FLOORS-1; furthest_floor > e.Elevators[e.Self_id].Floor; furthest_floor-- {
							if e.Elevators[e.Self_id].Internal_orders[1][furthest_floor] == 1 {
								break
							}
						}
					}
					Current_floor = ElevGetFloorSensorSignal()
					if Current_floor != -1 {
						e.Elevators[e.Self_id].Floor = Current_floor
						if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
							if e.Elevators[e.Self_id].Internal_orders[1][Current_floor] == 1 {
								if furthest_floor == Current_floor{
									e.Elevators[e.Self_id].Current_Dir = DIR_STOP
									e.StopAndOpenDoor(Current_floor)
									
									break
								}
							}
						} else if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
							if e.Elevators[e.Self_id].Internal_orders[0][Current_floor] == 1 {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.StopAndOpenDoor(Current_floor)
								
							}
						}
					}
				case DIR_DOWN:
					ElevSetMotorDirection(DIR_DOWN)
					
					if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
						for furthest_floor = 0; furthest_floor < e.Elevators[e.Self_id].Floor; furthest_floor++ {
							if e.Elevators[e.Self_id].Internal_orders[0][furthest_floor] == 1 {
								break
							}
						}
					}
					Current_floor = ElevGetFloorSensorSignal()
					if Current_floor != -1 {
						e.Elevators[e.Self_id].Floor = Current_floor
						if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
							if e.Elevators[e.Self_id].Internal_orders[0][Current_floor] == 1 {
								if furthest_floor == Current_floor{
									e.Elevators[e.Self_id].Current_Dir = DIR_STOP
									e.StopAndOpenDoor(Current_floor)
									break
								}
							}
						} else if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
							if e.Elevators[e.Self_id].Internal_orders[1][Current_floor] == 1 {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.StopAndOpenDoor(Current_floor)
								
							}
						}
					}

				case DIR_STOP:
					ElevSetMotorDirection(DIR_STOP)
					furthest_floor = -1
					e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
					for floor := 0; floor < N_FLOORS; floor++ {
						if e.Elevators[e.Self_id].Internal_orders[0][floor] == 1 {
							e.Elevators[e.Self_id].Planned_Dir = DIR_UP
							if e.Elevators[e.Self_id].Floor < floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_UP
								break
							} else if e.Elevators[e.Self_id].Floor > floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
								break
							}
						} else if e.Elevators[e.Self_id].Internal_orders[1][floor] == 1 {
							e.Elevators[e.Self_id].Planned_Dir = DIR_DOWN
								if e.Elevators[e.Self_id].Floor < floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_UP
								break
							} else if e.Elevators[e.Self_id].Floor > floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
								break
							}
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

	e.Elevators[e.Self_id].Internal_orders[0][floor] = 0
	e.Elevators[e.Self_id].Internal_orders[1][floor] = 0
	if e.Elevators[e.Self_id].Current_Dir == DIR_DOWN {
		e.Elevators[e.Self_id].External_orders[1][floor] = 1
	} else if e.Elevators[e.Self_id].Current_Dir == DIR_UP {
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

func calculateCost(elev Elevator, buttonType int, floor int) int {

	cost := int(math.Abs(float64((10*(floor - elev.Floor)))))
	numOrders := 0
	for x := 0; x < 4; x++{
		for y := 0; y < 3; y++{
			if elev.Internal_orders[y][x] == 1 {
				numOrders = numOrders + 1 
			}
		}
	}
	cost = cost + numOrders * 10

	if ((floor - elev.Floor) > 0) && (elev.Planned_Dir == DIR_UP) {
		if buttonType == BTN_UP {
			cost = cost - 100
		} 
	} else if (floor - elev.Floor) < 0 && elev.Planned_Dir == DIR_DOWN{
		if buttonType == BTN_DOWN {
			cost = cost - 100
		}
	} else if (floor - elev.Floor) == 0 {

		if elev.Planned_Dir == DIR_STOP {
			cost = cost - 80 
		}
	}
	Print("Kosten er: ")
	Print(cost)
	return cost

}


func costForRealsThough(elev Elevator, buttonType int, floor int) {

	if elev.Planned_Dir == -1 {
		//is the queue empty?
		elev.isQueueEmpty()	
	}


}
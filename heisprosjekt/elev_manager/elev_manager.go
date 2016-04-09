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
	Fsm_initiateElev()
	Println("ETTER")
	e.Elevators[e.Self_id].Current_Dir = DIR_STOP
	e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
	e.Elevators[e.Self_id].Current_Floor = 0
	e.Elevators[e.Self_id].State = STATE_IDLE
	e.Elevators[e.Self_id].Self_id = e.Self_id
	e.Elevators[e.Self_id].Furthest_Floor = -1
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

func (e *elev_manager) Em_processElevOrders(LampChan chan Message) {
	current_floor := -1	

	for {
		switch e.Elevators[e.Self_id].State {
		case STATE_IDLE:
			e.Elevators[e.Self_id].Current_Dir = DIR_STOP
			e.Elevators[e.Self_id].State = STATE_RUNNING
			/*for floor := 0; floor < 4; floor++ {
				for buttonType := 0; buttonType <= 3; buttonType++ {
					if e.Elevators[e.Self_id].Internal_orders[buttonType][floor] == 1 {
						e.Elevators[e.Self_id].State = STATE_RUNNING
					}
				}
			}*/
		case STATE_RUNNING:
				switch e.Elevators[e.Self_id].Current_Dir {
				case DIR_UP:
					ElevSetMotorDirection(DIR_UP)
				
					current_floor = ElevGetFloorSensorSignal()
					if current_floor != -1 {
						e.Elevators[e.Self_id].Current_Floor = current_floor
						if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][current_floor] == 1 {
							e.Elevators[e.Self_id].Current_Dir = DIR_STOP
							e.StopAndOpenDoor(current_floor, BTN_CMD, LampChan)
						}
						if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
							e.Elevators[e.Self_id].Furthest_Floor = e.calcFurthestFloor(e.Elevators[e.Self_id].Planned_Dir)

							if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][current_floor] == 1 {
								if e.Elevators[e.Self_id].Furthest_Floor == current_floor{
									e.Elevators[e.Self_id].Current_Dir = DIR_STOP
									e.StopAndOpenDoor(current_floor, BTN_DOWN, LampChan)
									break
								}
							}
						} else if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
							if e.Elevators[e.Self_id].Internal_orders[BTN_UP][current_floor] == 1 {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.StopAndOpenDoor(current_floor, BTN_UP, LampChan)
							}
						}
					}
				case DIR_DOWN:
					ElevSetMotorDirection(DIR_DOWN)
				
					current_floor = ElevGetFloorSensorSignal()
					if current_floor != -1 {
						e.Elevators[e.Self_id].Current_Floor = current_floor
						if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][current_floor] == 1 {
							e.Elevators[e.Self_id].Current_Dir = DIR_STOP
							e.StopAndOpenDoor(current_floor, BTN_CMD, LampChan)
						}
						if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
							e.Elevators[e.Self_id].Furthest_Floor = e.calcFurthestFloor(e.Elevators[e.Self_id].Planned_Dir)

							if e.Elevators[e.Self_id].Internal_orders[BTN_UP][current_floor] == 1 {
								if e.Elevators[e.Self_id].Furthest_Floor == current_floor{
									e.Elevators[e.Self_id].Current_Dir = DIR_STOP
									e.StopAndOpenDoor(current_floor, BTN_UP, LampChan)
									break
								}
							}
						} else if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
							if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][current_floor] == 1 {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.StopAndOpenDoor(current_floor, BTN_DOWN, LampChan)
							}
						}
					}
				case DIR_STOP:
					ElevSetMotorDirection(DIR_STOP)

					e.Elevators[e.Self_id].Furthest_Floor = -1
					e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
					e.check4OrdersAndDirChange(LampChan)
				}

		case STATE_DOOROPEN:
			ElevSetMotorDirection(DIR_STOP)
		case STATE_STOP:
			Println("Stop button pressed")
			ElevSetMotorDirection(DIR_STOP)
			//Set stop lamp
		}
	}
}

func (e *elev_manager) check4OrdersAndDirChange(LampChan chan Message) {
	for floor := 0; floor < N_FLOORS; floor++ {
		if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] == 1 {
			if e.Elevators[e.Self_id].Current_Floor < floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_UP
				e.Elevators[e.Self_id].Planned_Dir = DIR_UP
				break
			} else if e.Elevators[e.Self_id].Current_Floor > floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
				e.Elevators[e.Self_id].Planned_Dir = DIR_DOWN
				break
			}
		} else if e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] == 1 {
			e.Elevators[e.Self_id].Planned_Dir = DIR_UP
			if e.Elevators[e.Self_id].Current_Floor < floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_UP
				break
			} else if e.Elevators[e.Self_id].Current_Floor > floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
				break
			} else{
				e.StopAndOpenDoor(e.Elevators[e.Self_id].Current_Floor, BTN_UP, LampChan)
			}
		} else if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] == 1 {
			e.Elevators[e.Self_id].Planned_Dir = DIR_DOWN
			if e.Elevators[e.Self_id].Current_Floor < floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_UP
				break
			} else if e.Elevators[e.Self_id].Current_Floor > floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
				break
			} else{
				e.StopAndOpenDoor(e.Elevators[e.Self_id].Current_Floor, BTN_DOWN, LampChan)
			}
		}
	}
}

func (e *elev_manager) calcFurthestFloor(planned_dir int) int {
	switch planned_dir {
	case DIR_UP:
		for furthest_floor := 0; furthest_floor < e.Elevators[e.Self_id].Current_Floor; furthest_floor++ {
			if e.Elevators[e.Self_id].Internal_orders[BTN_UP][furthest_floor] == 1 {
				return furthest_floor
			}
		}
	case DIR_DOWN:
		for furthest_floor := N_FLOORS-1; furthest_floor > e.Elevators[e.Self_id].Current_Floor; furthest_floor-- {
			if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][furthest_floor] == 1 {
				return furthest_floor
			}
		}
	case DIR_STOP:
		Println("THIS IS THE WRONG PLANNED DIR BRUH")
		return -1
	}
	return e.Elevators[e.Self_id].Furthest_Floor
}


func (e *elev_manager) doorTimeout(floor int, button int) {
	ElevSetDoorOpenLamp(0)
	e.Em_RemoveOrders(floor, button)

	e.Elevators[e.Self_id].State = STATE_RUNNING
}

func (e *elev_manager) Em_RemoveOrders(floor int, button int) { //HEI BEDRE NAVN DA
	e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] = 0
	switch button {
	case BTN_UP:
		e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 0
	case BTN_DOWN:
		e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 0
	case BTN_CMD:
		if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
			e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 0
		} else if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
			e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 0
		}
	}
}

func (e *elev_manager) Em_AddInternalOrders(floor int, button int) {
	switch button {
	case BTN_UP:
		e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 1
	case BTN_DOWN:
		e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 1
	case BTN_CMD:
		e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] = 1
	}
}

func (e *elev_manager) Em_AddExternalOrders(floor int, buttonType int) {
	if buttonType == BTN_UP {
			e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 1
			e.Elevators[e.Self_id].External_orders[BTN_UP][floor] = 1

	} else if buttonType == BTN_DOWN {
			e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 1
			e.Elevators[e.Self_id].External_orders[BTN_DOWN][floor] = 1
		
	}
}


func updateFloorLight(floor int) {
	ElevSetFloorLamp(floor)
}

func (e *elev_manager) Em_NewElevator(elevMessage Message) {

	//e.Elevators[elevMessage.Source] = Elevator{}

}

func (e *elev_manager) StopAndOpenDoor(floor int, button int, lampChan chan Message) { //needs a better name
	ElevSetMotorDirection(DIR_STOP)
	Println("stopping")
	ElevSetDoorOpenLamp(1)
	updateButtonLamp(button, floor, lampChan)
	time.AfterFunc(time.Second*3, func() { e.doorTimeout(floor, button) })
	e.Elevators[e.Self_id].State = STATE_DOOROPEN
}

func updateButtonLamp(button int, floor int, lampChan chan Message) {
	switch button {
	case BTN_CMD:
		ElevSetButtonLamp(button, floor, 0)
	default:
		lampMessage := Message{ID: LampID, ButtonType: button, Floor: floor}
		lampChan <- lampMessage
	}
}

func (e *elev_manager) NewElevator(msg Message) {

}

func (e *elev_manager) RemoveElevator(target int) {
	if e.Self_id == target {

	}
}

func calculateCost(elev Elevator, buttonType int, floor int) int {

	cost := int(math.Abs(float64((10*(floor - elev.Current_Floor)))))
	numOrders := 0
	for x := 0; x < 4; x++{
		for y := 0; y < 3; y++{
			if elev.Internal_orders[y][x] == 1 {
				numOrders = numOrders + 1 
			}
		}
	}
	cost = cost + numOrders * 10

	if ((floor - elev.Current_Floor) > 0) && (elev.Planned_Dir == DIR_UP) {
		if buttonType == BTN_UP {
			cost = cost - 100
		} 
	} else if (floor - elev.Current_Floor) < 0 && elev.Planned_Dir == DIR_DOWN{
		if buttonType == BTN_DOWN {
			cost = cost - 100
		}
	} else if (floor - elev.Current_Floor) == 0 {

		if elev.Planned_Dir == DIR_STOP {
			cost = cost - 80 
		}
	}
	Print("Kosten er: ")
	Print(cost)
	return cost

}


func (e *elev_manager) determine_target_elev(button int, floor int) int{
	min := N_FLOORS
	var cost int
	ideal_elev := e.Self_id
	for key, elev := range e.Elevators {
		cost = 0
		
		if elev.Current_Floor == floor && (elev.State == STATE_IDLE || elev.State == STATE_DOOROPEN) {
			return key
		}

		if elev.Internal_orders[button][floor] == 1 {
			return key
		}

		cost = e.get_cost_for_order(key, elev, button, floor)

		if cost < min {
			ideal_elev = key
			min = cost
		}
	}
	return ideal_elev
}


func (e *elev_manager) get_cost_for_order(key int, elev *Elevator, button int, order_floor int) int {
	search_dir := elev.Planned_Dir
	search_floor := elev.Current_Floor

	order_floor_passed_once := false
	cost := 0

	for {
		if search_floor == 0 && search_dir == DIR_DOWN || search_floor == N_FLOORS-1 && search_dir == DIR_UP {
			search_dir = 0 - search_dir
		} else {
			search_floor += search_dir
		}

		if search_floor == order_floor {
			order_floor_passed_once = true
			if button == BTN_CMD || search_dir == DIR_UP && button == BTN_UP || search_dir == DIR_DOWN && button == BTN_DOWN {
				break
			}
		}
		
		if e.orders_on_floor_in_dir(key, search_dir, search_floor) {
			cost += 3
		} 
		
		if (e.more_orders_in_dir(key, search_dir, search_floor)) || (bool(math.Abs(float64(order_floor - elev.Current_Floor)) > math.Abs(float64(order_floor - search_floor)) && !order_floor_passed_once)) {
			cost += 1		}
	}

	return cost
}

func (e *elev_manager) more_orders_in_dir(id int, dir int, floor int) bool {
	for i := (floor + dir); i >= 0 && i < N_FLOORS; i = i + dir {
		if e.Elevators[id].Internal_orders[BTN_CMD][i] == 1 || e.Elevators[id].Internal_orders[BTN_UP][i] == 1 || e.Elevators[id].Internal_orders[BTN_DOWN][i] == 1 {
			return true
		}
	}
	return false
}

func (e *elev_manager) orders_on_floor_in_dir(id int, dir int, floor int) bool {
	//Vår heis

	if e.Elevators[id].Internal_orders[BTN_CMD][floor] == 1 || dir == DIR_UP && e.Elevators[id].Internal_orders[BTN_UP][floor] == 1 || dir == DIR_DOWN && e.Elevators[id].Internal_orders[BTN_DOWN][floor] == 1 {
		return true
	}
	return false
}

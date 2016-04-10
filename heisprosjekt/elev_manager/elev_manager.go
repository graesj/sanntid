package elev_manager

import (
	. ".././message"
	. ".././network"
	. ".././structs"
	. "./fsm"
	. "./fsm/driver"
	. "fmt"
	"math"
	"time"
)

type elev_manager struct {
	master          int
	OperationError  bool
	Self_id         int
	external_orders [2][N_FLOORS]int  //This is where the orders from the floor panels are put. These orders are broadcasted to all Elevators,
	Elevators       map[int]*Elevator //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func Em_makeElevManager() elev_manager {
	var e elev_manager
	e.Elevators = make(map[int]*Elevator)
	e.Self_id = GetLastNumbersOfIp()
	e.Elevators[e.Self_id] = new(Elevator)

	Fsm_initiateElev()
	e.OperationError = false
	e.Elevators[e.Self_id].Current_Dir = DIR_STOP
	e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
	e.Elevators[e.Self_id].Current_Floor = 0
	e.Elevators[e.Self_id].State = STATE_IDLE
	e.Elevators[e.Self_id].Self_id = e.Self_id
	e.Elevators[e.Self_id].Furthest_Floor = -1
	e.master = e.Self_id

	return e
}

func (e *elev_manager) Em_newElevator(elev Elevator) {

	e.Elevators[elev.Self_id] = &elev
	Print("Ny heis er lagt til med self_id:")
	Println(elev.Self_id)
	Println("Lagt inn i lista: ", e.Elevators[elev.Self_id].Self_id)
	if elev.Self_id < e.Self_id {
		e.master = elev.Self_id
		Print("Den nye heisen ble master")
	}
}

func (e *elev_manager) Em_isMaster() bool {
	if e.master == e.Self_id {
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
	engineTrouble := false
	lastFloor := e.Elevators[e.Self_id].Current_Floor
	engineCheck := time.NewTimer(3 * time.Second)

	for {

		switch e.Elevators[e.Self_id].State {
		case STATE_IDLE:
			select {
			case <-engineCheck.C:
			default:
			}
			e.Elevators[e.Self_id].Current_Dir = DIR_STOP
			e.Elevators[e.Self_id].Planned_Dir = DIR_STOP

			for floor := 0; floor < N_FLOORS; floor++ {
				for buttonType := 0; buttonType < 3; buttonType++ {
					if e.Elevators[e.Self_id].Internal_orders[buttonType][floor] == 1 {
						engineCheck.Reset(3 * time.Second)
						e.Elevators[e.Self_id].State = STATE_RUNNING
					}
				}
			}
		case STATE_RUNNING:

			select {
			case <-engineCheck.C:

				if lastFloor == e.Elevators[e.Self_id].Current_Floor {
					engineTrouble = true
					Println("hbu")
				} else {
					lastFloor = e.Elevators[e.Self_id].Current_Floor
					Println("WJA")
					engineCheck.Reset(2 * time.Second)
				}

			default:
				//asdfsd
			}

			if engineTrouble {
				e.OperationError = true
				engineTrouble = false
				break
			}

			switch e.Elevators[e.Self_id].Current_Dir {
			case DIR_UP:

				ElevSetMotorDirection(DIR_UP)
				current_floor = ElevGetFloorSensorSignal()
				if current_floor != -1 {
					updateFloorLight(current_floor)

					e.Elevators[e.Self_id].Current_Floor = current_floor
					if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][current_floor] == 1 {
						e.Elevators[e.Self_id].Current_Dir = DIR_STOP

						e.Elevators[e.Self_id].State = STATE_DOOROPEN
						break
					}
					if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
						e.Elevators[e.Self_id].Furthest_Floor = e.calcFurthestFloor(e.Elevators[e.Self_id].Planned_Dir)

						if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][current_floor] == 1 {
							if e.Elevators[e.Self_id].Furthest_Floor == current_floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.Elevators[e.Self_id].State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
						if e.Elevators[e.Self_id].Internal_orders[BTN_UP][current_floor] == 1 {
							e.Elevators[e.Self_id].Current_Dir = DIR_STOP
							e.Elevators[e.Self_id].State = STATE_DOOROPEN
							break
						}
					}
				}
			case DIR_DOWN:
				ElevSetMotorDirection(DIR_DOWN)

				current_floor = ElevGetFloorSensorSignal()
				if current_floor != -1 {
					updateFloorLight(current_floor)
					e.Elevators[e.Self_id].Current_Floor = current_floor
					if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][current_floor] == 1 {
						e.Elevators[e.Self_id].Current_Dir = DIR_STOP
						e.Elevators[e.Self_id].State = STATE_DOOROPEN
						break
					}
					if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
						e.Elevators[e.Self_id].Furthest_Floor = e.calcFurthestFloor(e.Elevators[e.Self_id].Planned_Dir)

						if e.Elevators[e.Self_id].Internal_orders[BTN_UP][current_floor] == 1 {
							if e.Elevators[e.Self_id].Furthest_Floor == current_floor {
								e.Elevators[e.Self_id].Current_Dir = DIR_STOP
								e.Elevators[e.Self_id].State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
						if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][current_floor] == 1 {
							e.Elevators[e.Self_id].Current_Dir = DIR_STOP
							e.Elevators[e.Self_id].State = STATE_DOOROPEN
							break
						}
					}
				}
			case DIR_STOP:
				ElevSetMotorDirection(DIR_STOP)
				boll := true
				for floor := 0; floor < N_FLOORS; floor++ {
					for buttonType := 0; buttonType < 3; buttonType++ {
						if e.Elevators[e.Self_id].Internal_orders[buttonType][floor] == 1 {
							boll = false
						}
					}
				}
				if boll {
					e.Elevators[e.Self_id].State = STATE_IDLE
					break
				}

				e.Elevators[e.Self_id].Furthest_Floor = -1
				e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
				e.check4OrdersAndDirChange(LampChan)
			}

		case STATE_DOOROPEN:
			ElevSetMotorDirection(DIR_STOP)
			ElevSetDoorOpenLamp(1)

			doorTimeout := time.NewTimer(3 * time.Second)
			<-doorTimeout.C
			//Kode etter
			ElevSetDoorOpenLamp(0)
			e.Em_RemoveOrders(current_floor, BTN_CMD, LampChan)
			e.Elevators[e.Self_id].State = STATE_RUNNING

			doorTimeout.Stop()
			engineCheck.Reset(3 * time.Second)
			ElevSetMotorDirection(DIR_STOP)
		}
	}
}

func (e *elev_manager) check4OrdersAndDirChange(LampChan chan Message) {
	foundFloor := 0
	for floor := 0; floor < N_FLOORS; floor++ {
		if e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] == 1 {

			if e.Elevators[e.Self_id].Current_Floor < floor {

				e.Elevators[e.Self_id].Current_Dir = DIR_UP
				e.Elevators[e.Self_id].Planned_Dir = DIR_UP
				foundFloor = 1
				break
			} else if e.Elevators[e.Self_id].Current_Floor > floor {
				e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
				e.Elevators[e.Self_id].Planned_Dir = DIR_DOWN
				foundFloor = 1
				break
			} else {
				e.StopAndOpenDoor(e.Elevators[e.Self_id].Current_Floor, BTN_CMD, LampChan)
				break
			}
		}
	}
	if foundFloor == 0 {
		for floor := 0; floor < N_FLOORS; floor++ {
			if e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] == 1 {
				e.Elevators[e.Self_id].Planned_Dir = DIR_UP
				if e.Elevators[e.Self_id].Current_Floor < floor {
					e.Elevators[e.Self_id].Current_Dir = DIR_UP
					break
				} else if e.Elevators[e.Self_id].Current_Floor > floor {
					e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
					break
				} else {
					e.StopAndOpenDoor(e.Elevators[e.Self_id].Current_Floor, BTN_UP, LampChan)
					break
				}
			} else if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] == 1 {
				e.Elevators[e.Self_id].Planned_Dir = DIR_DOWN
				if e.Elevators[e.Self_id].Current_Floor < floor {
					e.Elevators[e.Self_id].Current_Dir = DIR_UP
					break
				} else if e.Elevators[e.Self_id].Current_Floor > floor {
					e.Elevators[e.Self_id].Current_Dir = DIR_DOWN
					break
				} else {
					e.StopAndOpenDoor(e.Elevators[e.Self_id].Current_Floor, BTN_DOWN, LampChan)
					break
				}
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
		for furthest_floor := N_FLOORS - 1; furthest_floor > e.Elevators[e.Self_id].Current_Floor; furthest_floor-- {
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

func (e *elev_manager) doorTimeout(floor int, button int, lampChan chan Message) {
	ElevSetDoorOpenLamp(0)
	e.Em_RemoveOrders(floor, button, lampChan)

	e.Elevators[e.Self_id].State = STATE_RUNNING
}

func (e *elev_manager) Em_RemoveOrders(floor int, button int, lampChan chan Message) { //HEI BEDRE NAVN DA
	isEmpty := false
	e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] = 0
	ElevSetButtonLamp(button, floor, 0)
	if e.Elevators[e.Self_id].Planned_Dir == DIR_UP {
		isEmpty = e.shouldIChangeDir(DIR_UP)
		if isEmpty {
			e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 0
			e.Elevators[e.Self_id].External_orders[BTN_DOWN][floor] = 0

			updateButtonLamp(BTN_DOWN, floor, lampChan)
		}
		e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 0
		e.Elevators[e.Self_id].External_orders[BTN_UP][floor] = 0
		if floor != N_FLOORS-1 {
			updateButtonLamp(BTN_UP, floor, lampChan)
		}
	} else if e.Elevators[e.Self_id].Planned_Dir == DIR_DOWN {
		isEmpty = e.shouldIChangeDir(DIR_DOWN)
		if isEmpty {
			e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 0
			e.Elevators[e.Self_id].External_orders[BTN_UP][floor] = 0
			updateButtonLamp(BTN_UP, floor, lampChan)
		}
		e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 0
		e.Elevators[e.Self_id].External_orders[BTN_DOWN][floor] = 0
		if floor != 0 {
			updateButtonLamp(BTN_DOWN, floor, lampChan)
		}
	}
}

func (e *elev_manager) shouldIChangeDir(planned_dir int) bool {
	switch planned_dir {
	case DIR_UP:
		for floor := e.Elevators[e.Self_id].Current_Floor; floor < N_FLOORS; floor++ {
			if e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] == 1 || e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] == 1 {
				return false
			}
		}
	case DIR_DOWN:
		for floor := e.Elevators[e.Self_id].Current_Floor; floor >= 0; floor-- {
			if e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] == 1 || e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] == 1 {
				return false
			}
		}
	}
	return true
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
	ElevSetDoorOpenLamp(1)
	//time.AfterFunc(time.Second*3, func() { e.doorTimeout(floor, button, lampChan) })

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

func (e *elev_manager) Determine_target_elev(button int, floor int) int {
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
	search_dir := 1
	if elev.Planned_Dir != DIR_STOP {
		search_dir = elev.Planned_Dir
	}

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

		if (e.more_orders_in_dir(key, search_dir, search_floor)) || (bool(math.Abs(float64(order_floor-elev.Current_Floor)) > math.Abs(float64(order_floor-search_floor)) && !order_floor_passed_once)) {
			cost += 1
		}
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

func (e *elev_manager) ConnectionTimeout(id int, fromMain chan Message) {
	Println("NOEN HAR DISKONNECKTA")
	e.updateMaster(id)

	if id != e.Self_id && e.Em_isMaster() {
		Println("Dette skal bare skje med master")
		for i := 0; i < N_FLOORS; i++ {
			if e.Elevators[id].External_orders[0][i] == 1 {
				buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: BTN_UP, Floor: i}
				fromMain <- buttonMessage
				Print("OPP I ")
				Println(i + 1)
			} else if e.Elevators[id].External_orders[1][i] == 1 {
				buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: BTN_DOWN, Floor: i}
				fromMain <- buttonMessage
				Print("NED I ")
				Println(i + 1)
			}
		}

	}
	Print("HER GÅR DET IKKE GALT JINX")
	if id != e.Self_id {
		delete(e.Elevators, id)
	}

}

func (e *elev_manager) updateMaster(id int) {

	nyMaster := 1000

	for key, _ := range e.Elevators {
		Print(key)
		if key != id && key < nyMaster {
			nyMaster = key

		}
	}

	if nyMaster == 1000 {
		Println("Klare ikke finne ny master")
	} else {

		e.master = nyMaster

	}

}

func (e *elev_manager) CheckIfOrderIsTaken(message Message, fromMain chan Message) {

	_, present := e.Elevators[message.Target]

	if present {

		if e.Elevators[message.Target].External_orders[message.ButtonType][message.Floor] != 1 {
			fromMain <- message
			Println("Resending message")
		}

	}
}

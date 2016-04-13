package elev_manager

import (
	. "fmt"
	"math"


	. ".././message"
	. ".././network"
	. ".././structs"
	. "./fsm"
	. "./fsm/driver"
)


type elev_manager struct {
	master    int
	Self_id   int
	Elevators map[int]*Elevator //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
}

func MakeElevManager() elev_manager {
	var e elev_manager
	e.Elevators = make(map[int]*Elevator)
	e.Self_id = GetLastNumbersOfIp()
	e.Elevators[e.Self_id] = new(Elevator)

	RunToFirstFloor()
	e.Elevators[e.Self_id].Current_Dir = DIR_STOP
	e.Elevators[e.Self_id].Planned_Dir = DIR_STOP
	e.Elevators[e.Self_id].ErrorType = ERROR_NONE
	e.Elevators[e.Self_id].Current_Floor = 0
	e.Elevators[e.Self_id].State = STATE_IDLE
	e.Elevators[e.Self_id].Self_id = e.Self_id
	e.Elevators[e.Self_id].Furthest_Floor = -1
	e.master = e.Self_id

	return e
}
func (e *elev_manager) NewElevator(elev Elevator) {

	e.Elevators[elev.Self_id] = &elev
	Print("Ny heis er lagt til med self_id:")
	Println(elev.Self_id)
	Println("Lagt inn i lista: ", e.Elevators[elev.Self_id].Self_id)
	if elev.Self_id < e.Self_id {
		e.master = elev.Self_id
		Print("Den nye heisen ble master")
	}
}

func (e *elev_manager) IsMaster() bool {
	if e.master == e.Self_id {
		return true
	} else {
		return false
	}
}

func (e *elev_manager) OnConnectionTimeout(id int, fromMain chan Message, elev Elevator) {
	Print("Disconnected : ")
	Println(id)
	e.UpdateMaster(id)

	if id != e.Self_id && e.IsMaster() {
		Println("Dette skal bare skje med master")
		for i := 0; i < N_FLOORS; i++ {
			if e.Elevators[id].External_orders[0][i] == 1 {
				buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: BTN_UP, Floor: i}
				fromMain <- buttonMessage
			} else if e.Elevators[id].External_orders[1][i] == 1 {
				buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: BTN_DOWN, Floor: i}
				fromMain <- buttonMessage
			}
		}

	}
	if id != e.Self_id {
		e.Elevators[id].ErrorType = elev.ErrorType
	}

}



/*
func (e *elev_manager) AddInternalOrders(floor int, button int) {
	switch button {
	case BTN_UP:
		e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 1
	case BTN_DOWN:
		e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 1
	case BTN_CMD:
		ElevSetButtonLamp(BTN_CMD, floor, 1)
		e.Elevators[e.Self_id].Internal_orders[BTN_CMD][floor] = 1
	}
}*/

/*
func (e *elev_manager) AddExternalOrders(floor int, buttonType int) {
	if buttonType == BTN_UP {
		e.Elevators[e.Self_id].Internal_orders[BTN_UP][floor] = 1
		e.Elevators[e.Self_id].External_orders[BTN_UP][floor] = 1

	} else if buttonType == BTN_DOWN {
		e.Elevators[e.Self_id].Internal_orders[BTN_DOWN][floor] = 1
		e.Elevators[e.Self_id].External_orders[BTN_DOWN][floor] = 1

	}
}*/

func (e *elev_manager) ElevatorUpdate(elev Elevator) {
	e.Elevators[elev.Self_id] = &elev
}

func (e *elev_manager) DetermineTargetElev(button int, floor int) int {
	min := N_FLOORS
	var cost int
	ideal_elev := e.Self_id
	for key, elev := range e.Elevators {

		if elev.ErrorType == ERROR_NONE {

			cost = 0

			if elev.Current_Floor == floor && (elev.State == STATE_IDLE || elev.State == STATE_DOOROPEN) {
				return key
			}

			if elev.Internal_orders[button][floor] == 1 {
				return key
			}

			cost = e.GetCostForOrder(key, elev, button, floor)

			if cost < min {
				ideal_elev = key
				min = cost
			}
		}
	}
	return ideal_elev
}

func (e *elev_manager) GetCostForOrder(key int, elev *Elevator, button int, order_floor int) int {
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

		if e.OrdersOnFloorInDir(key, search_dir, search_floor) {
			cost += 3
		}

		if (e.MoreOrdersInDir(key, search_dir, search_floor)) || (bool(math.Abs(float64(order_floor-elev.Current_Floor)) > math.Abs(float64(order_floor-search_floor)) && !order_floor_passed_once)) {
			cost += 1
		}
	}

	return cost
}

func (e *elev_manager) MoreOrdersInDir(id int, dir int, floor int) bool {

	for i := (floor + dir); i >= 0 && i < N_FLOORS; i = i + dir {

		if e.Elevators[id].Internal_orders[BTN_CMD][i] == 1 || e.Elevators[id].Internal_orders[BTN_UP][i] == 1 || e.Elevators[id].Internal_orders[BTN_DOWN][i] == 1 {
			return true

		}
	}

	return false
}

func (e *elev_manager) OrdersOnFloorInDir(id int, dir int, floor int) bool {
	//Vår heis
	if e.Elevators[id].Internal_orders[BTN_CMD][floor] == 1 || dir == DIR_UP && e.Elevators[id].Internal_orders[BTN_UP][floor] == 1 || dir == DIR_DOWN && e.Elevators[id].Internal_orders[BTN_DOWN][floor] == 1 {
		return true
	}
	return false
}

func (e *elev_manager) UpdateMaster(id int) {

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

func (e *elev_manager) CheckIfOrderIsReceived(message Message, fromMain chan Message) {

	_, present := e.Elevators[message.Target]

	if present {

		if e.Elevators[message.Target].External_orders[message.ButtonType][message.Floor] != 1 {
			fromMain <- message
			Println("Resending message")
		}

	}
}

func (e *elev_manager) ResendInitialOrders(fromMain chan Message, id int) {
	m := Message{ID: GET_UP_TO_DATE, Target: id, Elevator: *e.Elevators[id]}
	fromMain <- m
}

func (e *elev_manager) CopyInternalOrder(elev Elevator) {
	Println("Henter tilbake køen min")
	Println(elev.Internal_orders)
	for i := 0; i < N_FLOORS; i++ {
		if elev.Internal_orders[BTN_CMD][i] == 1 {
			Println(i)
			ElevSetButtonLamp(BTN_CMD, i, 1)
			AddInternalOrders(e.Elevators[e.Self_id],i, BTN_CMD)
		}
	}
}


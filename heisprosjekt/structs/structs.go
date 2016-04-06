package structs

type Elevator struct {
	State           int
	Dir             int
	Floor           int
	Self_id 		int 
	Internal_orders [2][N_FLOORS]byte //both external and internal orders
	External_orders [2][N_FLOORS]byte //orders from the external panel
	//Just for backup
}

const (
	N_FLOORS = 4
	DIR_UP   = 1
	DIR_DOWN = -1
	DIR_STOP = 0

	BTN_UP   = 0
	BTN_DOWN = 1
	BTN_CMD  = 2
)
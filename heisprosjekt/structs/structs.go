package structs

type Elevator struct {
	State           int
	Current_Dir     int
	Planned_Dir     int
	Current_Floor   int
	Furthest_Floor  int
	Self_id         int
	ErrorType       int
	Internal_orders [3][N_FLOORS]byte //both external and internal orders
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

	ERROR_NONE    = 0
	ERROR_NETWORK = 1
	ERROR_MOTOR   = 2
)

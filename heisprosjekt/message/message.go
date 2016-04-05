package message

const (
	IP          = 1
	ELEV_STRUCT = 2
	BUTTON      = 3
)

type Message struct {
	Source int
	Floor  int
	Target int
	e      Elevator
	Id     int
}

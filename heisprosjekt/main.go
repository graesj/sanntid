package main 

import "C"

import(
	"fmt"
	. "./network"
	. "./message"
	"time"
)

func main(){

	//var msg Message 
	themessage := Message{1,2,3}
	msg := make(chan Message,1)
	msg <- themessage
	msgrec := make(chan Message,1)


	go UDPsend(msg)
	go UDPListen(msgrec)
	fmt.Println("Hei") 

 for {
 	select {
 	case m:= <-msgrec:
 		fmt.Println(m)
 	case <- time.Tick(10 * time.Second):
 		msg <- Message{1,1,1}


 	}

 }
/*
io_set_bit(int channel)

func io_set_bit(channel int)
	C.io_set_bit(C.int(channel))
*/
}
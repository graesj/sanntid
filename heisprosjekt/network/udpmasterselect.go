package network_manager

import (
		. "UDP"
		"time"
		"net"
		".././message"
		"fmt"
)

func ip_broadcast(ip int, UDPsend chan Message){
	
	UDPsend <- Message{IP = ip, ID = 1};
	time.Sleep(100*Millisecond)

}


func network_manager(fromMain chan Message, toMain chan Message){


UDPsend := make(chan Message,50);
UDPrecieve := make(chan Message,50);

//


}




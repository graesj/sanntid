package main 

import(

	. "./elev_manager"
	"fmt"
)

func main(){
	for {
		fmt.Println(Em_checkButtons())
	}
}
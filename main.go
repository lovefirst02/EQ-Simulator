package main

import (
	"fmt"
	"simulator/Global"
	"simulator/Router"
	"simulator/Simulator"
)

var AsrsMap Simulator.ASRS

func main() {
	for i := 0; i < 3; i++ {
		AsrsID := fmt.Sprintf("Asrs%d", i+1)
		Asrs := Simulator.NewAsrs(AsrsID)
		Global.Asrs[AsrsID] = Asrs
		go Asrs.AsrsSimulator()
	}
	Router.InitRouter()
}

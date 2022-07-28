package Global

import (
	"simulator/Simulator"
)

var Asrs = make(map[string]*Simulator.ASRS)

var Lifter = make(map[string]*Simulator.LIFTER)

var Erack = make(map[string]*Simulator.ERACK)

var Port string

var EQcount int

var MCS string

var Mode string

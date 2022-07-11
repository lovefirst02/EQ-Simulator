package main

import (
	"fmt"
	"simulator/Global"
	"simulator/Router"
	"simulator/Simulator"

	"github.com/spf13/viper"
)

var AsrsMap Simulator.ASRS

func init() {
	viper.AddConfigPath("conf")
	viper.AddConfigPath(".")
	viper.SetConfigName("Setting")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("no such config file")
		} else {
			fmt.Println("read config error")
		}
		fmt.Println(err)
	}
	Global.Port = viper.GetString("Port")
	Global.EQcount = viper.GetInt("EQCOUNT")
	Global.MCS = viper.GetString("MCS")
	Global.Mode = viper.GetString("MODE")
}

func main() {
	for i := 0; i < Global.EQcount; i++ {
		AsrsID := fmt.Sprintf("Asrs%d", i+1)
		Asrs := Simulator.NewAsrs(AsrsID)
		Global.Asrs[AsrsID] = Asrs
		go Asrs.AsrsSimulator()
	}
	Router.InitRouter()
}

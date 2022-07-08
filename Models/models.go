package Models

import (
	"encoding/json"
	"time"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Control struct {
	AsrsID    string `json:"AsrsID"`
	Type      string `json:"Type"`
	MissionID string `json:"MissionID"`
}

type Mission struct {
	MissionID  string `json:"MissionID"`
	Sourceport string `json:"Sourceport"`
	Destport   string `json:"Destport"`
	CarrierID  string `json:"CarrierID"`
	Priority   string `json:"Priority"`
	Quantity   int    `json:"Quantity"`
	AsrsID     string `json:"AsrsID"`
}

type MissionStatus struct {
	MissionID  string `json:"MissionID"`
	Sourceport string `json:"Sourceport"`
	Destport   string `json:"Destport"`
	CarrierID  string `json:"CarrierID"`
	Priority   string `json:"Priority"`
	Quantity   int    `json:"Quantity"`
	AsrsID     string `json:"AsrsID"`
	Status     int    `json:"Status"`
}

type ASRS struct {
	AsrsID string    `json:"AsrsID"`
	Type   string    `json:"Type"`
	Status string    `json:"Status"`
	Time   time.Time `json:"Time"`
}

type AsrsAlarm struct {
	AsrsID string `json:"AsrsID"`
	ALID   int    `json:"ALID"`
	ALMSG  string `json:"ALMSG"`
}

func (asrs *ASRS) MarshalJSON() ([]byte, error) {
	type Alias ASRS
	return json.Marshal(&struct {
		*Alias
		Time string `json:"Time"`
	}{
		Alias: (*Alias)(asrs),
		Time:  asrs.Time.Format("2006-01-02 15:04:05"),
	})
}

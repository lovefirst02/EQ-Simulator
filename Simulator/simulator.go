package Simulator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"simulator/Models"
	"simulator/Util"
	"sync"
	"time"
)

type ASRS struct {
	AsrsID  string    `json:"AsrsID"`
	Type    string    `json:"Type"`
	Status  string    `json:"Status"`
	Time    time.Time `json:"Time"`
	Mission []Models.Mission
	Control chan Models.Control
	mux     sync.Mutex
}

func (asrs *ASRS) request_alarm_mcs() {
	alarm_slice := []string{
		"Alarm1",
		"Alarm2",
		"Alarm3",
	}
	rand.Seed(time.Now().Unix())
	random_alarm := alarm_slice[rand.Intn(len(alarm_slice))]
	alarm_data := Models.AsrsAlarm{
		AsrsID: asrs.AsrsID,
		ALID:   9999,
		ALMSG:  random_alarm,
	}
	client := &http.Client{}
	alarm_json, _ := json.Marshal(alarm_data)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/api/rs/device/alarm", bytes.NewBuffer(alarm_json))
	req.SetBasicAuth("admin", "motorcon")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println(alarm_data.AsrsID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println(err)
			return
		}
		// fmt.Println(string(body))
	}

}

func (asrs *ASRS) request_mcs(mission_status Models.MissionStatus) {
	client := &http.Client{}
	stauts_JSON, _ := json.Marshal(mission_status)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/api/rs/device/mission/status", bytes.NewBuffer(stauts_JSON))
	req.SetBasicAuth("admin", "motorcon")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println(mission_status.MissionID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println(err)
			return
		}
		// fmt.Println(string(body))
	}
	// fmt.Println(mission_status)
}

func (asrs *ASRS) asrsmissionsimulator(mission Models.Mission) {
	asrs.Status = "RUN"
	asrs.mux.Lock()
	mission_status := Models.MissionStatus{
		MissionID:  mission.MissionID,
		Sourceport: mission.Sourceport,
		Destport:   mission.Destport,
		CarrierID:  mission.CarrierID,
		Priority:   mission.CarrierID,
		Quantity:   mission.Quantity,
		AsrsID:     mission.AsrsID,
	}
	for {
		rand.Seed(time.Now().UnixNano())
		randomNum := Util.Random(1, 5)
		switch randomNum {
		case 1:
			mission_status.Status = 1
		case 2:
			mission_status.Status = 2
		case 3:
			mission_status.Status = 3
		case 4:
			mission_status.Status = 4
		}
		asrs.request_mcs(mission_status)
		if mission_status.Status == 3 {
			asrs.Status = "IDLE"
			asrs.mux.Unlock()
			return
		}
		time.Sleep(10 * time.Second)
	}
}

// func (asrs *ASRS) Init() {
// 	asrs.Control = make(chan Models.Control)
// 	asrs.Mission = make([]Models.Mission, 0)
// }

func (asrs *ASRS) DeleteMission(mission string) {
	go func() {
		i := 0
		for _, v := range asrs.Mission {
			if v.MissionID != mission {
				asrs.Mission[i] = v
				i++
			}
		}
		asrs.Mission = asrs.Mission[:i]
	}()
}

func (asrs *ASRS) AsrsControl(command Models.Control) {
	go func() {

		asrs.Control <- command
	}()
}

func (asrs *ASRS) AsrsMission(mission Models.Mission) {
	asrs.Mission = append(asrs.Mission, mission)
}

func (asrs *ASRS) AsrsSimulator() {
	// asrs.Init()
	for {
		if len(asrs.Mission) > 0 && !Util.MutexLocked(&asrs.mux) {
			mission := asrs.Mission[0]
			asrs.Mission = asrs.Mission[1:]

			go asrs.asrsmissionsimulator(mission)
		}
		if asrs.Status != "RUN" {
			rand.Seed(time.Now().UnixNano())
			randomNum := Util.Random(1, 3)
			select {
			case command := <-asrs.Control:
				switch command.Type {
				case "EMO":
					asrs.Status = "MAINTAIN"
				case "RESTART":
					if asrs.Status == "MAINTAIN" {
						asrs.Status = "IDLE"
					}
				case "CANCEL":
					asrs.DeleteMission(command.MissionID)
				case "DELETE":
					asrs.DeleteMission(command.MissionID)
				}
			default:
				if asrs.Status != "MAINTAIN" {
					switch randomNum {
					case 1:
						asrs.Status = "IDLE"
					case 2:
						asrs.request_alarm_mcs()
						asrs.Status = "ALARM"
					}
				}

			}
		}
		asrs.Time = time.Now()
		time.Sleep(5 * time.Second)
	}
}

func NewAsrs(AsrsID string) *ASRS {
	return &ASRS{
		AsrsID:  AsrsID,
		Mission: make([]Models.Mission, 0),
		Control: make(chan Models.Control),
	}
}

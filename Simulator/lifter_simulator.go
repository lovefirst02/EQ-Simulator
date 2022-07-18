package Simulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"simulator/Models"
	"simulator/Util"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type LIFTER struct {
	LifterID string    `json:"AsrsID"`
	Type     string    `json:"Type"`
	Status   string    `json:"Status"`
	Time     time.Time `json:"Time"`
	Mission  []Models.LifterMission
	Control  chan Models.LiterControl
	mux      sync.Mutex
}

func (lifter *LIFTER) request_alarm_mcs() {
	alarm_slice := []string{
		"Alarm1",
		"Alarm2",
		"Alarm3",
	}
	rand.Seed(time.Now().Unix())
	random_alarm := alarm_slice[rand.Intn(len(alarm_slice))]
	alarm_data := Models.AsrsAlarm{
		AsrsID: lifter.LifterID,
		ALID:   9999,
		ALMSG:  random_alarm,
	}
	client := &http.Client{}
	alarm_json, _ := json.Marshal(alarm_data)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/rs/device/alarm", viper.GetString("MCS")), bytes.NewBuffer(alarm_json))
	req.SetBasicAuth("admin", "motorcon")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
	}

}

func (lifter *LIFTER) request_mcs(mission_status Models.LifterMission) {
	client := &http.Client{}
	stauts_JSON, _ := json.Marshal(mission_status)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/rs/device/mission/status", viper.GetString("MCS")), bytes.NewBuffer(stauts_JSON))
	req.SetBasicAuth("admin", "motorcon")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
	}
}

func (lifter *LIFTER) liftermissionsimulator(mission Models.LifterMission) {
	lifter.Status = "RUN"
	lifter.mux.Lock()
	alarm_count := 0
	for {
		rand.Seed(time.Now().UnixNano())
		randomNum := Util.Random(1, 5)
		select {
		case msg := <-mission.Control:
			if msg == "RESOLVE" {
				mission.Status = 2
			}
		default:
			if mission.Status != 4 || lifter.Status == "RUN" {
				switch randomNum {
				case 1:
					mission.Status = 1
				case 2:
					mission.Status = 2
				case 3:
					mission.Status = 3
				case 4:
					if alarm_count <= 2 {
						mission.Status = 4
						lifter.Type = fmt.Sprintf("%s - ALARM", mission.MissionID)
						lifter.Status = "ALARM"
						lifter.request_alarm_mcs()
						alarm_count += 1
					}
				}
			}
			fmt.Printf("%s,%d\n", mission.MissionID, mission.Status)
			lifter.request_mcs(mission)
			if mission.Status == 3 {
				lifter.Status = "IDLE"
				fmt.Printf("%s,Complete\n", mission.MissionID)
				lifter.mux.Unlock()
				return
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (lifter *LIFTER) DeleteMission(mission string) {
	go func() {
		i := 0
		for _, v := range lifter.Mission {
			if v.MissionID != mission {
				lifter.Mission[i] = v
				i++
			}
		}
		lifter.Mission = lifter.Mission[:i]
	}()
}

func (lifter *LIFTER) LifterControl(command Models.LiterControl) {
	go func() {

		lifter.Control <- command
	}()
}

func (lifter *LIFTER) LifterMissionPrivateControl(Control Models.LifterMissionPrivateControl) bool {
	for _, v := range lifter.Mission {
		if v.MissionID == Control.MissionID {
			v.Control <- Control.Type
			return true
		}
	}
	return false
}

func (lifter *LIFTER) LifterMission(mission Models.LifterMission) {
	lifter.Mission = append(lifter.Mission, mission)
}

func (lifter *LIFTER) LifterSimulator() {
	// asrs.Init()
	for {
		if len(lifter.Mission) > 0 && !Util.MutexLocked(&lifter.mux) {
			mission := lifter.Mission[0]
			lifter.Mission = lifter.Mission[1:]
			go lifter.liftermissionsimulator(mission)
		}
		if lifter.Status != "RUN" {
			rand.Seed(time.Now().UnixNano())
			// randomNum := Util.Random(1, 3)
			select {
			case command := <-lifter.Control:
				switch command.Type {
				case "EMO":
					lifter.Status = "MAINTAIN"
				case "RESTART":
					if lifter.Status == "MAINTAIN" {
						lifter.Status = "IDLE"
					}
				case "CANCEL":
					lifter.DeleteMission(command.MissionID)
				case "DELETE":
					lifter.DeleteMission(command.MissionID)
				case "BUG_ALARM":
					lifter.Status = "BUG_ALARM"
					lifter.Type = fmt.Sprintf("%s - ALARM", command.MissionID)
				case "RESOLVE":
					if lifter.Status == "ALARM" {
						lifter.Type = ""
						lifter.Status = "RUN"
					} else {
						lifter.Type = ""
						lifter.Status = "RESOLVE_ALARM"
					}
				}
			default:
				if !(lifter.Status == "MAINTAIN" || lifter.Status == "ALARM") {
					lifter.Status = "IDLE"
					// switch randomNum {
					// case 1:
					// lifter.Status = "IDLE"
					// case 2:
					// lifter.Status = "IDLE"
					// lifter.request_alarm_mcs()
					// lifter.Status = "ALARM"
					// }
				}

			}
		}
		time.Sleep(5 * time.Second)
	}
}

func NewLifter(LifterID string) *LIFTER {
	return &LIFTER{
		LifterID: LifterID,
		Mission:  make([]Models.LifterMission, 0),
		Control:  make(chan Models.LiterControl),
	}
}

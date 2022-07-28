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
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
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

func (asrs *ASRS) request_mcs(mission_status Models.Mission) {
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

func (asrs *ASRS) request_erack(Device, Storage, Carrier string) {
	client := &http.Client{}
	data := Models.ErackStorage{
		ErackID:   Device,
		CarrierID: Carrier,
		Storage:   Storage,
	}
	json_data, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/device/api/erack/install", "http://127.0.0.1:8880"), bytes.NewBuffer(json_data))
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

func (asrs *ASRS) report_mcs(Device, DeviceLocation, CarrierID string, Event int) {
	client := &http.Client{}
	data := Models.EVENT{
		Device:          Device,
		Device_Location: DeviceLocation,
		CarrierID:       CarrierID,
		Event:           Event,
	}
	json_data, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/eq/report", viper.GetString("MCS")), bytes.NewBuffer(json_data))
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

func (asrs *ASRS) request_status_mcs(ID, status, Type string, time time.Time) {
	client := &http.Client{}
	status_model := &Models.AsrsStatus{
		AsrsID: ID,
		Type:   Type,
		Status: status,
		Time:   time,
	}
	stauts_JSON, _ := json.Marshal(status_model)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/rs/device/status", viper.GetString("MCS")), bytes.NewBuffer(stauts_JSON))
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

func (asrs *ASRS) asrsmissionsimulator(mission Models.Mission) {
	asrs.Status = "RUN"
	asrs.mux.Lock()
	alarm_count := 0
	run_count := 0
	max_random := 5
	mission.Status = 1
	asrs.report_mcs(asrs.AsrsID, asrs.AsrsID, mission.CarrierID, mission.Status)
	for {
		if alarm_count >= 2 {
			max_random = 4
		}
		rand.Seed(time.Now().UnixNano())
		randomNum := Util.Random(1, max_random)
		select {
		case msg := <-mission.Control:
			if msg == "RESOLVE" {
				mission.Status = 2
			}
		default:
			if mission.Status != 4 || asrs.Status == "RUN" {
				switch randomNum {
				case 1:
					// mission.Status = 1
				case 2:
					mission.Status = 2
					run_count = run_count + 1
				case 3:
					if run_count >= 2 {
						mission.Status = 3
					}
				case 4:
					if alarm_count <= 2 {
						mission.Status = 4
						asrs.Type = fmt.Sprintf("%s - ALARM", mission.MissionID)
						asrs.request_alarm_mcs()
						asrs.Status = "ALARM"
						alarm_count += 1
					}
				}
			}
			fmt.Printf("%s,%s,%s,%d\n", mission.MissionID, asrs.AsrsID, asrs.Status, mission.Status)
			asrs.request_mcs(mission)
			asrs.report_mcs(asrs.AsrsID, asrs.AsrsID, mission.CarrierID, mission.Status)
			if mission.Status == 3 && run_count >= 2 {
				if strings.Contains(mission.Destport, "ERACK") {
					DestPort_split := strings.Split(mission.Destport, "-")
					asrs.request_erack(DestPort_split[0], fmt.Sprintf("%s-%s", DestPort_split[1], DestPort_split[2]), mission.CarrierID)
				}
				asrs.Status = "IDLE"
				fmt.Printf("%s,Complete\n", mission.MissionID)
				asrs.mux.Unlock()
				return
			}
		}
		time.Sleep(5 * time.Second)
	}
}

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

func (asrs *ASRS) AsrsMissionPrivateControl(Control Models.MissionPrivateControl) bool {
	for _, v := range asrs.Mission {
		if v.MissionID == Control.MissionID {
			v.Control <- Control.Type
			return true
		}
	}
	return false
}

func (asrs *ASRS) AsrsMission(mission Models.Mission) {
	fmt.Println(mission)
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
			// randomNum := Util.Random(1, 3)
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
				case "BUG_ALARM":
					asrs.Status = "BUG_ALARM"
					asrs.Type = fmt.Sprintf("%s - ALARM", command.MissionID)
				case "RESOLVE":
					if asrs.Status == "ALARM" {
						asrs.Type = ""
						asrs.Status = "RUN"
					} else {
						asrs.Type = ""
						asrs.Status = "RESOLVE_ALARM"
					}
				}
			default:
				if !(asrs.Status == "MAINTAIN" || asrs.Status == "ALARM") {
					asrs.Status = "IDLE"
				}

			}
		}
		asrs.Time = time.Now()
		asrs.request_status_mcs(asrs.AsrsID, asrs.Status, asrs.Type, asrs.Time)
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

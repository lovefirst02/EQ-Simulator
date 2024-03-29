package Simulator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"simulator/Models"
	"time"

	"github.com/spf13/viper"
)

const (
	Empty StorageStatus = iota
	Install
	Uninstall
	Pre
)

type ERACK struct {
	ErackID string
	Storage map[string]StorageDetail
}

type StorageDetail struct {
	CarrierID   string
	InstallTime string
	RemoveTime  string
	Status      StorageStatus
}

type StorageStatus int

func (erack *ERACK) report_mcs(Device, DeviceLocation, CarrierID string, Event int) {
	client := &http.Client{
		Timeout: 10 * time.Millisecond,
	}
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

func (erack *ERACK) init(col int, row int) {
	erack.Storage = make(map[string]StorageDetail, row*col)
	for i := 0; i < col; i++ {
		for x := 0; x < row; x++ {
			storage := fmt.Sprintf("%d-%d", i+1, x+1)
			storagedetail := &StorageDetail{
				Status: Uninstall,
			}
			erack.Storage[storage] = *storagedetail
		}
	}
}

func (erack *ERACK) Install(Storage, CarrierID string) error {
	if data, ok := erack.Storage[Storage]; ok {
		if data.Status != Install {
			data.CarrierID = CarrierID
			data.InstallTime = time.Now().Format("2006-01-02 15:04:05")
			data.RemoveTime = ""
			data.Status = Install
			erack.Storage[Storage] = data
			erack.report_mcs(erack.ErackID, Storage, data.CarrierID, int(Install))
		}
	} else {
		return errors.New("no storage")
	}
	return nil
}

func (erack *ERACK) Uninstall(Storage string) error {
	if data, ok := erack.Storage[Storage]; ok {
		if data.Status == Install {
			data.CarrierID = ""
			data.RemoveTime = time.Now().Format("2006-01-02 15:04:05")
			data.Status = Uninstall
			erack.Storage[Storage] = data
			erack.report_mcs(erack.ErackID, Storage, data.CarrierID, int(Uninstall))
		}
	} else {
		return errors.New("no storage")
	}
	return nil
}

func (erack *ERACK) PreStorage(Storage string) error {
	if data, ok := erack.Storage[Storage]; ok {
		if data.Status != Pre && data.Status != Install {
			data.CarrierID = "Pre"
			data.InstallTime = time.Now().Format("2006-01-02 15:04:05")
			data.Status = Pre
			erack.Storage[Storage] = data
			erack.report_mcs(erack.ErackID, Storage, data.CarrierID, int(Pre))
		}
	} else {
		return errors.New("no storage")
	}
	return nil
}

func (erack *ERACK) CheckEmptyStorage() []string {
	result := make([]string, 0, len(erack.Storage))
	for k, v := range erack.Storage {
		if !(v.Status == Install || v.Status == Pre) {
			result = append(result, k)
		}
	}
	return result
}

func (erack *ERACK) ErackSimulator() {
	for {
		time.Sleep(5 * time.Second)
	}
}

func NewErack(ErackID string, col int, row int) *ERACK {
	Erack := &ERACK{
		ErackID: ErackID,
	}
	Erack.init(col, row)
	return Erack
}

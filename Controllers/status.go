package Controllers

import (
	"regexp"
	"simulator/Global"
	"simulator/Models"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	var status []Models.AsrsStatus

	for _, v := range Global.Asrs {
		status_data := Models.AsrsStatus{
			AsrsID: v.AsrsID,
			Type:   v.Type,
			Status: v.Status,
			Time:   v.Time,
		}
		status = append(status, status_data)
	}

	sort.SliceStable(status, func(i, j int) bool {
		re := regexp.MustCompile("[0-9]+")
		id1 := re.FindAllString(status[i].AsrsID, -1)[0]
		id2 := re.FindAllString(status[j].AsrsID, -1)[0]
		id1_int, _ := strconv.Atoi(id1)
		id2_int, _ := strconv.Atoi(id2)
		return id1_int < id2_int
	})

	c.JSON(200, gin.H{
		"Code":        0,
		"Description": status,
	})
}

////////////////////////////Status////////////////////////////////

func LifterStatus(c *gin.Context) {
	var status []Models.LifterStatus

	for _, v := range Global.Lifter {
		status_data := Models.LifterStatus{
			LifterID: v.LifterID,
			Type:     v.Type,
			Status:   v.Status,
			Time:     v.Time,
		}
		status = append(status, status_data)
	}

	sort.SliceStable(status, func(i, j int) bool {
		re := regexp.MustCompile("[0-9]+")
		id1 := re.FindAllString(status[i].LifterID, -1)[0]
		id2 := re.FindAllString(status[j].LifterID, -1)[0]
		id1_int, _ := strconv.Atoi(id1)
		id2_int, _ := strconv.Atoi(id2)
		return id1_int < id2_int
	})

	c.JSON(200, gin.H{
		"Code":        0,
		"Description": status,
	})
}

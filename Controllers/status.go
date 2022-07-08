package Controllers

import (
	"simulator/Global"
	"simulator/Models"
	"sort"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	var status []Models.ASRS

	for _, v := range Global.Asrs {
		status_data := Models.ASRS{
			AsrsID: v.AsrsID,
			Type:   v.Type,
			Status: v.Status,
			Time:   v.Time,
		}
		status = append(status, status_data)
	}

	sort.SliceStable(status, func(i, j int) bool {
		return status[i].AsrsID[len(status[i].AsrsID)-1] < status[j].AsrsID[len(status[j].AsrsID)-1]
	})

	c.JSON(200, gin.H{
		"Code":        0,
		"Description": status,
	})
}

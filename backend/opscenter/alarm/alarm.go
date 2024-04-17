package alarm

import (
	"encoding/json"
	"kubehostwarden/db"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AlarmInfo struct {
	Id         string  `json:"id" gorm:"primaryKey"`
	UserId     string  `json:"userId" gorm:"column:user_id"`
	HostId     string  `json:"hostId" gorm:"column:host_id"`
	AlarmType  string  `json:"alarmType" gorm:"column:alarm_type"`
	AlarmValue float64 `json:"alarmValue" gorm:"column:alarm_value"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (AlarmInfo) TableName() string {
	return "alarm_info"
}

type AlarmInfoRequest struct {
	UserId     string  `json:"userId" gorm:"column:user_id"`
	HostId     string  `json:"hostId" gorm:"column:host_id"`
	AlarmType  string  `json:"alarmType" gorm:"column:alarm_type"`
	AlarmValue float64 `json:"alarmValue" gorm:"column:alarm_value"`
}

func SetAlarm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var alarmInfoReq AlarmInfoRequest
	err := json.NewDecoder(r.Body).Decode(&alarmInfoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}

	alarmInfo := AlarmInfo{
		Id:         uuid.NewString()[:8],
		UserId:     alarmInfoReq.UserId,
		HostId:     alarmInfoReq.HostId,
		AlarmType:  alarmInfoReq.AlarmType,
		AlarmValue: alarmInfoReq.AlarmValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.GetMysqlClient().Client.Save(&alarmInfo)

	w.Write([]byte(`successfully set alarm`))
}

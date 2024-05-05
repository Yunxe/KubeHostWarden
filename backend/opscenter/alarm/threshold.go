package alarm

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/opscenter/logger"
	"kubehostwarden/types"
	"kubehostwarden/utils/constant"
	resp "kubehostwarden/utils/responsor"
	"kubehostwarden/utils/scheduler"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type setThresholdReq struct {
	HostId    string                 `json:"host_id" validate:"required"`
	Metric    string                 `json:"metric" validate:"required"`
	SubMetric string                 `json:"sub_metric" validate:"required"`
	Threshold float64                `json:"threshold" validate:"required"`
	Type      constant.ThresholdType `json:"type" validate:"required"`
}

func SetThreshold(ctx context.Context, req setThresholdReq) resp.Responsor {
	var user types.User
	userId := ctx.Value(constant.UserIDKey).(string)
	res := db.GetMysqlClient().Client.WithContext(ctx).Where("id = ?", userId).First(&user)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get user: %v", res.Error),
		}
	}

	var threshold = &types.ThresholdInfo{
		Id:        uuid.NewString(),
		HostId:    req.HostId,
		Metric:    req.Metric,
		SubMetric: req.SubMetric,
		Threshold: req.Threshold,
		Type:      string(req.Type),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryId, err := scheduler.AddJob("@every 10s", func() {
		CheckThreshold(ctx, user.Email, req.HostId, req.Metric, req.SubMetric, req.Threshold, req.Type)
	})
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to add job: %v", err),
		}
	}

	threshold.EntryId = int(entryId)
	res = db.GetMysqlClient().Client.WithContext(ctx).Save(threshold)
	if res.Error != nil {
		scheduler.RemoveJob(cron.EntryID(threshold.EntryId))
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to save threshold: %v", res.Error),
		}
	}

	logger.Info(userId, "阈值设置成功", "阈值:", threshold)
	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "threshold set successfully",
		Result:  threshold,
	}
}

func GetThreshold(ctx context.Context, values url.Values) resp.Responsor {
	hostId := values.Get("hostId")
	var thresholds []types.ThresholdInfo
	res := db.GetMysqlClient().Client.WithContext(ctx).Where("host_id = ?", hostId).Find(&thresholds)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get thresholds: %v", res.Error),
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "thresholds fetched successfully",
		Result:  thresholds,
	}
}

type deleteThresholdReq struct {
	Id string `json:"id" validate:"required"`
}

func DeleteThreshold(ctx context.Context, req deleteThresholdReq) resp.Responsor {
	var threshold = &types.ThresholdInfo{
		Id: req.Id,
	}
	res := db.GetMysqlClient().Client.WithContext(ctx).First(threshold)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get threshold: %v", res.Error),
		}
	}
	scheduler.RemoveJob(cron.EntryID(threshold.EntryId))

	res = db.GetMysqlClient().Client.WithContext(ctx).Delete(threshold)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete threshold: %v", res.Error),
		}
	}

	logger.Info(ctx.Value(constant.UserIDKey).(string), "阈值删除成功", "阈值:", threshold)
	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "threshold deleted successfully",
	}
}

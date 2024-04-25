package host

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	mysql "kubehostwarden/db"
	"kubehostwarden/types"
	"kubehostwarden/utils/constant"
	"kubehostwarden/utils/logger"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"net/url"
)

type probeHelper struct {
	sshInfo types.SSHInfo
	host    *types.Host
}

func Register(ctx context.Context, sshInfo types.SSHInfo) resp.Responsor {
	var pHelper probeHelper

	pHelper.sshInfo = sshInfo

	err := pHelper.probe(ctx)
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to probe host: %v", err),
		}
	}

	var owner types.User
	ownerId := ctx.Value(constant.UserIDKey).(string)
	result := mysql.GetMysqlClient().Client.WithContext(ctx).Where("id = ?", ownerId).First(&owner)
	if result.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to query host owner: %v", result.Error),
		}
	}

	pHelper.host.OwnerId = owner.Id
	pHelper.host.Owner = owner.Username

	result = mysql.GetMysqlClient().Client.WithContext(ctx).Create(pHelper.host)
	if result.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to save host: %v", result.Error),
		}
	}

	err = pHelper.createPod(ctx)
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to create pod: %v", err),
		}
	}
	logger.Info("pod created successfully", "host", pHelper.host.Id)

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "host registered successfully",
		Result:  pHelper.host,
	}
}

func Retrieve(ctx context.Context, values url.Values) resp.Responsor {
	var hosts []types.Host
	ownerId := ctx.Value(constant.UserIDKey).(string)
	result := db.GetMysqlClient().Client.WithContext(ctx).Where("owner_id = ?", ownerId).Find(&hosts)
	if result.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to query host: %v", result.Error),
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "host retrieved successfully",
		Result:  hosts,
	}
}

type DeleteReq struct {
	HostId string `json:"hostId"`
}

func Delete(ctx context.Context, req DeleteReq) resp.Responsor {
	result := mysql.GetMysqlClient().Client.WithContext(ctx).Delete(&types.Host{Id: req.HostId})
	if result.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete host: %v", result.Error),
		}
	}

	//delete pod
	err := deletePod(ctx, "host-"+req.HostId)
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete pod: %v", err),
		}
	}

	//TODO need transaction

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "host deleted successfully",
	}
}

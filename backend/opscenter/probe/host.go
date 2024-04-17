package probe

import (
	"context"
	"fmt"
	mysql "kubehostwarden/db"
	"kubehostwarden/types"
	"kubehostwarden/utils/constant"
	resp "kubehostwarden/utils/responsor"
	"net/http"

	"golang.org/x/crypto/ssh"
)

type probeHelper struct {
	sshInfo   types.SSHInfo
	sshClient *ssh.Client
	host      *types.Host
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

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "host registered successfully",
		Result:  pHelper.host,
	}
}

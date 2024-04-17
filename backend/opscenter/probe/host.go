package probe

import (
	"context"
	"fmt"
	mysql "kubehostwarden/db"
	"kubehostwarden/types"
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

	result := mysql.GetMysqlClient().Client.WithContext(ctx).Create(pHelper.host)
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

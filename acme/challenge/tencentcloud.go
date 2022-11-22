package challenge

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"github.com/paraparty/acme-task/model"
)

func TencentCloudChallenge(client *lego.Client, task *model.Task) error {
	tencentCloudConfig := tencentcloud.NewDefaultConfig()
	tencentCloudConfig.SecretID = task.Challenge.Credential.SecretID
	tencentCloudConfig.SecretKey = task.Challenge.Credential.SecretKey
	tencentCloudProvider, err := tencentcloud.NewDNSProviderConfig(tencentCloudConfig)
	if err != nil {
		return err
	}

	err = client.Challenge.SetDNS01Provider(tencentCloudProvider)
	if err != nil {
		return err
	}

	return nil
}

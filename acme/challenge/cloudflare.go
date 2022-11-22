package challenge

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/paraparty/acme-task/model"
)

func CloudflareChallenge(client *lego.Client, task *model.Task) error {
	cloudflareConfig := cloudflare.NewDefaultConfig()
	cloudflareConfig.ZoneToken = task.Challenge.Credential.ZoneToken
	cloudflareConfig.AuthToken = task.Challenge.Credential.AuthToken
	cloudflareDnsProvider, err := cloudflare.NewDNSProviderConfig(cloudflareConfig)
	if err != nil {
		return err
	}

	err = client.Challenge.SetDNS01Provider(cloudflareDnsProvider)
	if err != nil {
		return err
	}

	return nil
}

package acme

import (
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/paraparty/acme-task/model"
)

func NewClient(user *model.User) (*lego.Client, error) {
	acmeConfig := lego.NewConfig(user)
	acmeConfig.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	acmeConfig.Certificate.KeyType = certcrypto.EC256
	acmeConfig.Certificate.Timeout = 7 * 24 * time.Hour

	client, err := lego.NewClient(acmeConfig)
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	return client, nil
}

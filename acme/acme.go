package acme

import (
	"context"
	"encoding/json"
	"fmt"

	publicca "cloud.google.com/go/security/publicca/apiv1beta1"
	"cloud.google.com/go/security/publicca/apiv1beta1/publiccapb"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/paraparty/acme-task/model"
	"google.golang.org/api/option"
)

func NewClient(config *model.Config, user *model.User) (*lego.Client, error) {
	acmeConfig := lego.NewConfig(user)
	acmeConfig.Certificate.KeyType = certcrypto.EC256
	if config.Acme.Type == "google" {
		acmeConfig.CADirURL = "https://dv.acme-v02.api.pki.goog/directory"
		err := getGoogleAcmeAuth(config)
		if err != nil {
			return nil, err
		}
	} else if config.Acme.Type == "test" {
		acmeConfig.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else if config.Acme.Type == "r3" {
		acmeConfig.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	} else {
		return nil, fmt.Errorf("acme type invalid")
	}

	client, err := lego.NewClient(acmeConfig)
	if err != nil {
		return nil, err
	}

	account := registration.RegisterEABOptions{
		TermsOfServiceAgreed: true,
		HmacEncoded:          config.Acme.HmacEncoded,
		Kid:                  config.Acme.KeyId,
	}

	reg, err := client.Registration.RegisterWithExternalAccountBinding(account)
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	return client, nil
}

func getGoogleAcmeAuth(config *model.Config) error {
	credential, err := json.Marshal(config.Acme.Details.Credential)
	if err != nil {
		return err
	}

	ctx := context.Background()
	c, err := publicca.NewPublicCertificateAuthorityClient(ctx, option.WithCredentialsJSON(credential))
	if err != nil {
		return err
	}
	defer c.Close()

	req := &publiccapb.CreateExternalAccountKeyRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", config.Acme.Details.Project),
		ExternalAccountKey: &publiccapb.ExternalAccountKey{},
	}
	resp, err := c.CreateExternalAccountKey(ctx, req)
	if err != nil {
		return err
	}
	// TODO: Use resp.
	config.Acme.KeyId = resp.KeyId
	config.Acme.HmacEncoded = string(resp.B64MacKey)
	return nil
}

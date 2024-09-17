package challenge

import (
	"context"
	"fmt"

	"github.com/go-acme/lego/v4/lego"
	"github.com/paraparty/acme-task/imagex"
	"github.com/paraparty/acme-task/model"
	volc "github.com/volcengine/volc-sdk-golang/base"
	volcImageX "github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

func ImageXChallenge(client *lego.Client, task *model.Task) error {
	imagexService := imagex.CreateImageXService()

	imagexService.SetCredential(volc.Credentials{
		AccessKeyID:     task.Challenge.Credential.SecretID,
		SecretAccessKey: task.Challenge.Credential.SecretKey,
	})

	imagexChallenge := &imagexChallenge{
		Service: imagexService,
	}
	err := imagexChallenge.init()
	if err != nil {
		return nil
	}

	err = client.Challenge.SetHTTP01Provider(imagexChallenge)
	if err != nil {
		return err
	}

	return nil
}

type imagexChallenge struct {
	Service       *volcImageX.Imagex
	DomainMapping map[string]string
}

func ChallengePath(token string) string {
	return ".well-known/acme-challenge/" + token
}

func (c *imagexChallenge) Present(domain, token, keyAuth string) error {
	path := ChallengePath(token)
	serviceId, ok := c.DomainMapping[domain]
	if !ok {
		return fmt.Errorf("%s not found in ImageX", domain)
	}

	arg := &volcImageX.ApplyUploadImageParam{
		ServiceId: serviceId,
		UploadNum: 1,
		StoreKeys: []string{path},
	}
	_, err := c.Service.UploadImages(arg, [][]byte{[]byte(keyAuth)})
	if err != nil {
		return err
	}

	return nil
}

func (c *imagexChallenge) CleanUp(domain, token, keyAuth string) error {
	path := ChallengePath(token)
	serviceId, ok := c.DomainMapping[domain]
	if !ok {
		return fmt.Errorf("%s not found in ImageX", domain)
	}

	req := volcImageX.DeleteImageUploadFilesReq{
		DeleteImageUploadFilesQuery: &volcImageX.DeleteImageUploadFilesQuery{
			ServiceID: serviceId,
		},
		DeleteImageUploadFilesBody: &volcImageX.DeleteImageUploadFilesBody{
			StoreUris: []string{path},
		},
	}

	_, err := c.Service.DeleteImageUploadFiles(context.Background(), &req)
	if err != nil {
		return err
	}

	return nil
}

func (c *imagexChallenge) init() error {
	getAllImageServicesQuery := &volcImageX.GetAllImageServicesQuery{}
	servicesResp, err := c.Service.GetAllImageServices(context.Background(), getAllImageServicesQuery)
	if err != nil {
		return err
	}
	services := servicesResp.Result

	c.DomainMapping = make(map[string]string, 0)

	for _, item := range services.Services {
		for _, domain := range item.DomainInfos {
			c.DomainMapping[domain.DomainName] = item.ServiceID
		}
	}
	return nil
}

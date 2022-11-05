package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"github.com/paraparty/acme-task/acme"
	"github.com/paraparty/acme-task/configuration"
	"github.com/paraparty/acme-task/handler"
	"github.com/paraparty/acme-task/model"
)

func main() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	config, err := configuration.ReadConfig()

	user := &model.User{
		Email: config.Email,
		Key:   privateKey,
	}

	acmeClient, err := acme.NewClient(user)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range config.Tasks {
		err = doTask(acmeClient, &task)
		if err != nil {
			log.Printf("%v", err)
		}
	}

}

func doTask(client *lego.Client, task *model.Task) error {
	err := resolveChallenge(client, task)
	if err != nil {
		return err
	}

	request := certificate.ObtainRequest{
		Domains: task.Domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	if task.TaskDetails.Type == "imagex" {
		return handler.ImageXHandler(task, certificates)
	}
	return fmt.Errorf("task type not support")
}

func resolveChallenge(client *lego.Client, task *model.Task) error {
	if task.Challenge.Type == "tencent-cloud" {
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

	return fmt.Errorf("credential type not support")
}

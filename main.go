package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/paraparty/acme-task/acme"
	"github.com/paraparty/acme-task/acme/challenge"
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
	if err != nil {
		log.Fatal(err)
	}

	user := &model.User{
		Email: config.Acme.Email,
		Key:   privateKey,
	}

	acmeClient, err := acme.NewClient(config, user)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range config.Tasks {
		err = doTask(config, acmeClient, &task)
		if err != nil {
			log.Printf("%v", err)
		}
	}

}

func doTask(config *model.Config, client *lego.Client, task *model.Task) error {
	err := resolveChallenge(client, task)
	if err != nil {
		return err
	}

	request := certificate.ObtainRequest{
		Domains: task.Domains,
		Bundle:  true,
	}
	if config.Acme.Type == "google" {
		if config.Acme.ValidityPeriod != "" {
			duration, err := time.ParseDuration(config.Acme.ValidityPeriod)
			if err == nil {
				request.NotAfter = time.Now().Add(duration)
			} else {
				log.Printf("ValidityPeriod is not acceptable, %v", err)
			}
		}
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	if task.TaskDetails.Type == "imagex" {
		return handler.ImageXHandler(task, certificates)
	} else if task.TaskDetails.Type == "file" {
		return handler.CertFileHandler(task, certificates)
	}
	return fmt.Errorf("task type not support")
}

func resolveChallenge(client *lego.Client, task *model.Task) error {
	if task.Challenge.Type == "cloudflare" {
		return challenge.CloudflareChallenge(client, task)
	} else if task.Challenge.Type == "tencent-cloud" {
		return challenge.TencentCloudChallenge(client, task)
	} else if task.Challenge.Type == "imagex" {
		return challenge.ImageXChallenge(client, task)
	}

	return fmt.Errorf("credential type not support")
}

package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/paraparty/acme-task/imagex"
	"github.com/paraparty/acme-task/model"
	volc "github.com/volcengine/volc-sdk-golang/base"
)

func ImageXHandler(task *model.Task, certificates *certificate.Resource) error {
	if len(task.Domains) <= 0 || len(task.TaskDetails.Services) <= 0 {
		return fmt.Errorf("no need to run task")
	}

	imagexService := imagex.CreateImageXService()

	imagexService.SetCredential(volc.Credentials{
		AccessKeyID:     task.TaskDetails.Credential.SecretID,
		SecretAccessKey: task.TaskDetails.Credential.SecretKey,
	})

	addedCert, err := imagex.AddCert(imagexService, certificates)
	if err != nil {
		return err
	}

	servicesInfo, err := imagexService.GetImageServices("")
	if err != nil {
		return err
	}

	for _, service := range servicesInfo.Services {
		log.Printf("now processing service %s(%s)", service.ServiceName, service.ServiceId)

		if !arrContains(task.TaskDetails.Services, service.ServiceId) {
			log.Printf("skip service:%s(%s) for service id not hit", service.ServiceName, service.ServiceId)
			continue
		}

		for _, domain := range service.DomainInfos {
			if !checkDomain(task.Domains, domain.DomainName) {
				log.Printf("skip service:%s(%s) for domain not hit", service.ServiceName, service.ServiceId)
				continue
			}

			err := imagex.EnableServiceHttps(imagexService, service.ServiceId, domain.DomainName, addedCert.CertId)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

			log.Printf("set cert for %s(%s):%s finished", service.ServiceName, service.ServiceId, domain.DomainName)
		}
	}

	return nil
}

func checkDomain(domainsList []string, domain string) bool {
	if arrContains(domainsList, domain) {
		return true
	}

	domainStep := strings.Split(domain, ".")
	domainSub := strings.Join(domainStep[1:], ".")
	return arrContains(domainsList, "*."+domainSub)
}

func arrContains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

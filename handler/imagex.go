package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/paraparty/acme-task/imagex"
	"github.com/paraparty/acme-task/model"
	volc "github.com/volcengine/volc-sdk-golang/base"
	volcImagex "github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

func ImageXHandler(task *model.Task, certificates *certificate.Resource) error {
	ctx := context.Background()

	if len(task.Domains) <= 0 || len(task.TaskDetails.Services) <= 0 {
		return fmt.Errorf("no need to run task")
	}

	imagexService := imagex.CreateImageXService()

	imagexService.SetCredential(volc.Credentials{
		AccessKeyID:     task.TaskDetails.Credential.SecretID,
		SecretAccessKey: task.TaskDetails.Credential.SecretKey,
	})

	var addedCert *model.AddCertResponse
	err := retry.Do(func() error {
		var err error
		addedCert, err = imagex.AddCert(imagexService, certificates)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(time.Second*5), retry.OnRetry(func(n uint, err error) {
		log.Printf("add cert error: retry:%d err:%+v", n, err)
	}))
	if err != nil {
		return err
	}

	var servicesInfoResp *volcImagex.GetAllImageServicesRes
	getAllImageServicesQuery := &volcImagex.GetAllImageServicesQuery{}
	err = retry.Do(func() error {
		var err error
		servicesInfoResp, err = imagexService.GetAllImageServices(ctx, getAllImageServicesQuery)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.OnRetry(func(n uint, err error) {
		log.Printf("get service info error: retry:%d err:%+v", n, err)
	}))
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	servicesInfo := servicesInfoResp.Result

	for _, service := range servicesInfo.Services {
		log.Printf("now processing service %s(%s)", service.ServiceName, service.ServiceID)

		if !arrContains(task.TaskDetails.Services, service.ServiceID) {
			log.Printf("skip service:%s(%s) for service id not hit", service.ServiceName, service.ServiceID)
			continue
		}

		for _, domain := range service.DomainInfos {
			if !checkDomain(task.Domains, domain.DomainName) {
				log.Printf("skip service:%s(%s) for domain not hit", service.ServiceName, service.ServiceID)
				continue
			}

			setCertErr := retry.Do(func() error {
				retryErr := imagex.EnableServiceHttps(ctx, imagexService, service.ServiceID, domain.DomainName, addedCert.CertId)
				if retryErr != nil {
					return retryErr
				}
				return nil
			}, retry.Attempts(3), retry.Delay(time.Second*5), retry.OnRetry(func(n uint, err error) {
				log.Printf("set cert for %s(%s):%s retry:%d err:%+v", service.ServiceName, service.ServiceID, domain.DomainName, n, err)
			}))
			if setCertErr != nil {
				log.Printf("%v", setCertErr)
				continue
			}

			log.Printf("set cert for %s(%s):%s finished", service.ServiceName, service.ServiceID, domain.DomainName)
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

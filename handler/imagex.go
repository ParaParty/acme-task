package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/google/uuid"
	"github.com/paraparty/acme-task/model"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	volc "github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/imagex"
)

func ImageXHandler(task *model.Task, certificates *certificate.Resource) error {
	if len(task.Domains) <= 0 || len(task.TaskDetails.Services) <= 0 {
		return fmt.Errorf("no need to run task")
	}

	imagexService := createImageXService()

	imagexService.SetCredential(volc.Credentials{
		AccessKeyID:     task.TaskDetails.Credential.SecretID,
		SecretAccessKey: task.TaskDetails.Credential.SecretKey,
	})

	addedCert, err := addCert(imagexService, certificates)
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

			err := connectServiceDomain(imagexService, service.ServiceId, domain.DomainName, addedCert.CertId)
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

func addCert(instance *imagex.ImageX, certificates *certificate.Resource) (*model.AddCertResponse, error) {
	certSuffix, _ := uuid.NewUUID()
	req := &model.AddCertRequest{
		Name:    "auto-acme-task-" + time.Now().Format(time.RFC3339) + "-" + certSuffix.String(),
		Public:  string(certificates.Certificate),
		Private: string(certificates.PrivateKey),
	}

	resp := &model.AddCertResponse{}
	err := instance.ImageXPost("AddCert", url.Values{}, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func connectServiceDomain(instance *imagex.ImageX, serviceId string, domain string, certId string) error {
	query := url.Values{}
	query.Add("ServiceId", serviceId)

	req := &model.UpdateHttpsRequest{
		Domain: domain,
		Https: model.UpdateHttpsItemRequest{
			CertId:      certId,
			EnableHttp2: true,
			EnableHttps: true,
		},
	}

	resp := common.StringPtr("")

	err := instance.ImageXPost("UpdateHttps", query, req, resp)
	if err != nil {
		return err
	}
	return err
}

func arrContains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func createImageXService() *imagex.ImageX {
	instance := imagex.NewInstanceWithRegion(volc.RegionCnNorth1)
	instance.ApiInfoList["AddCert"] = &volc.ApiInfo{
		Method: http.MethodPost,
		Path:   "/",
		Query: url.Values{
			"Action":  []string{"AddCert"},
			"Version": []string{"2018-08-01"},
		},
	}
	instance.ApiInfoList["UpdateHttps"] = &volc.ApiInfo{
		Method: http.MethodPost,
		Path:   "/",
		Query: url.Values{
			"Action":  []string{"UpdateHttps"},
			"Version": []string{"2018-08-01"},
		},
	}
	return instance
}

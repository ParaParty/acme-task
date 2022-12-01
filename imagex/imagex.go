package imagex

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/google/uuid"
	"github.com/paraparty/acme-task/model"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	volc "github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/imagex"
)

func CreateImageXService() *imagex.ImageX {
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

func AddCert(instance *imagex.ImageX, certificates *certificate.Resource) (*model.AddCertResponse, error) {
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

func EnableServiceHttps(instance *imagex.ImageX, serviceId string, domain string, certId string) error {
	query := url.Values{}
	query.Add("ServiceId", serviceId)

	req := &model.UpdateHttpsRequest{
		Domain: domain,
		Https: model.UpdateHttpsItemRequest{
			CertId:              certId,
			EnableHttp2:         true,
			EnableHttps:         true,
			EnableForceRedirect: true,
			RedirectCode:        "301",
			ForceRedirectType:   "http2https",
			TlsVersions:         []string{"tlsv1.2", "tlsv1.3"},
		},
	}

	resp := common.StringPtr("")

	err := instance.ImageXPost("UpdateHttps", query, req, resp)
	if err != nil {
		return err
	}
	return err
}

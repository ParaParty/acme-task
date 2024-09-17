package imagex

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/google/uuid"
	"github.com/paraparty/acme-task/model"
	volc "github.com/volcengine/volc-sdk-golang/base"
	imagex "github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

func CreateImageXService() *imagex.Imagex {
	instance := imagex.NewInstanceWithRegion(volc.RegionCnNorth1)
	instance.ApiInfoList["AddCert"] = &volc.ApiInfo{
		Method: http.MethodPost,
		Path:   "/",
		Query: url.Values{
			"Action":  []string{"AddCert"},
			"Version": []string{"2018-08-01"},
		},
	}
	return instance
}

func AddCert(instance *imagex.Imagex, certificates *certificate.Resource) (*model.AddCertResponse, error) {
	certSuffix, _ := uuid.NewUUID()
	req := &model.AddCertRequest{
		Name:    "auto-main-" + time.Now().Format(time.RFC3339) + "-" + certSuffix.String(),
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

func EnableServiceHttps(ctx context.Context, instance *imagex.Imagex, serviceId string, domain string, certId string) error {
	req := &imagex.UpdateHTTPSReq{

		UpdateHTTPSQuery: &imagex.UpdateHTTPSQuery{
			ServiceID: serviceId,
		},
		UpdateHTTPSBody: &imagex.UpdateHTTPSBody{
			Domain: domain,
			HTTPS: &imagex.UpdateHTTPSBodyHTTPS{
				CertID:              certId,
				EnableHTTP2:         true,
				EnableHTTPS:         true,
				EnableForceRedirect: true,
				ForceRedirectCode:   "301",
				ForceRedirectType:   "http2https",
				TLSVersions:         []string{"tlsv1.2", "tlsv1.3"},
			},
		},
	}

	_, err := instance.UpdateHTTPS(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

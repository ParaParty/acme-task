package model

type UpdateHttpsRequest struct {
	Domain string                 `json:"domain"`
	Https  UpdateHttpsItemRequest `json:"https"`
}

type UpdateHttpsItemRequest struct {
	CertId      string `json:"cert_id"`
	EnableHttp2 bool   `json:"enable_http2"`
	EnableHttps bool   `json:"enable_https"`
}

type AddCertRequest struct {
	Name    string `json:"name"`
	Public  string `json:"public"`
	Private string `json:"private"`
}

type AddCertResponse struct {
	CertId     string `json:"cert_id"`
	CertName   string `json:"cert_name"`
	CommonName string `json:"common_name"`
	CreateTime int64  `json:"create_time"`
	NotAfter   int64  `json:"not_after"`
	Issuer     string `json:"issuer"`
}

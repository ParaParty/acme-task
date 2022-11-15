package handler

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/paraparty/acme-task/model"
	"gopkg.in/yaml.v3"
)

type secretYaml struct {
	Kind       string             `yaml:"kind"`
	ApiVersion string             `yaml:"apiVersion"`
	Metadata   secretMetadataYaml `yaml:"metadata"`
	Data       map[string]string  `yaml:"data"`
	Type       string             `yaml:"type"`
}

type secretMetadataYaml struct {
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations"`
}

func CertFileHandler(task *model.Task, certificates *certificate.Resource) error {
	if task.TaskDetails.CertificatePath != "" {
		f, err := os.Create(task.TaskDetails.CertificatePath)
		if err != nil {
			return err
		}
		_, err = f.Write(certificates.Certificate)
		if err != nil {
			return err
		}
	}

	if task.TaskDetails.PrivateKeyPath != "" {
		f, err := os.Create(task.TaskDetails.PrivateKeyPath)
		if err != nil {
			return err
		}
		_, err = f.Write(certificates.PrivateKey)
		if err != nil {
			return err
		}
	}

	if task.TaskDetails.OutputPath != "" {
		secretYaml := secretYaml{
			Kind:       "Secret",
			ApiVersion: "v1",
			Type:       "kubernetes.io/tls",
			Metadata: secretMetadataYaml{
				Name: fmt.Sprintf("tls-%s", certificates.Domain),
				Annotations: map[string]string{
					"cert-manager.io/alt-names":        certificates.Domain,
					"cert-manager.io/certificate-name": fmt.Sprintf("tls-%s", certificates.Domain),
					"cert-manager.io/common-name":      certificates.Domain,
					"cert-manager.io/ip-sans":          "",
					"cert-manager.io/issuer-group":     "cert-manager.io",
					"cert-manager.io/issuer-kind":      "ClusterIssuer",
					"cert-manager.io/issuer-name":      "letsencrypt",
					"cert-manager.io/uri-sans":         "",
				},
			},
			Data: map[string]string{
				"tls.crt": "",
				"tls.key": "",
			},
		}

		secretYaml.Data["tls.crt"] = base64.StdEncoding.EncodeToString(certificates.Certificate)
		secretYaml.Data["tls.key"] = base64.StdEncoding.EncodeToString(certificates.PrivateKey)

		result, err := yaml.Marshal(secretYaml)
		if err != nil {
			return err
		}
		f, err := os.Create(task.TaskDetails.OutputPath)
		if err != nil {
			return err
		}
		_, err = f.Write(result)
		if err != nil {
			return err
		}
	}

	return nil
}

package configuration

import (
	"encoding/json"
	"os"

	"github.com/paraparty/acme-task/model"
	"github.com/paraparty/acme-task/utils"
)

func ReadConfig() (*model.Config, error) {
	config := &model.Config{}
	var bytes []byte
	var err error

	err = os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", "true")
	if err != nil {
		return nil, err
	}

	acmeTaskConfigFilePath := utils.GetEnvVar("acme_task_config_file", "")
	acmeTaskConfigFile := utils.GetEnvVar("acme_task_config", "")
	if acmeTaskConfigFilePath != "" {
		bytes, err = os.ReadFile(acmeTaskConfigFilePath)
		if err != nil {
			return nil, err
		}
	} else if acmeTaskConfigFile != "" {
		bytes = []byte(acmeTaskConfigFile)
	} else {
		bytes, err = os.ReadFile("config.json")
		if err != nil {
			return nil, err
		}
	}

	err = json.Unmarshal(bytes, config)
	return config, err
}

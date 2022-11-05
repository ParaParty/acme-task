package configuration

import (
	"encoding/json"

	"github.com/paraparty/acme-task/model"
	"github.com/paraparty/acme-task/utils"
)

func ReadConfig() (*model.Config, error) {
	config := &model.Config{}
	err := json.Unmarshal([]byte(utils.GetEnvVar("acme-task-config", "{}")), config)
	return config, err
}

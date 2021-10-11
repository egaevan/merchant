package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/egaevan/merchant/constant"
	"github.com/egaevan/merchant/model"
)

func GetConfig() (*model.Config, error) {
	cfg := &model.Config{}

	jsonFile, err := ioutil.ReadFile(constant.ConfigProjectFilepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFile, &cfg)

	return cfg, nil
}

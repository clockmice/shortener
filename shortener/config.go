package shortener

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type Config struct {
	DbConfig       DBConfig `json:"db_config"`
	UrlAliasLength int      `json:"url_alias_length"`
	ServerPort     string   `json:"server_port"`
	UrlHost        string   `json:"url_host"`
}

type DBConfig struct {
	DbUsername string `json:"db_username"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	Table      string `json:"table"`
}

func ReadConfig(path string) (*Config, error) {
	var config Config

	r, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Could not read file '%v'. %v\n", path, err)
		return &config, err
	}

	err = json.Unmarshal(r, &config)
	if err != nil {
		err = fmt.Errorf("Could not unmarshall json. %v\n", err)
		return &config, err
	}

	return &config, nil
}

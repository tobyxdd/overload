package main

import (
	"encoding/json"
	"io/ioutil"
)

type interfaceConfig struct {
	Name   string `json:"name"`
	Weight uint   `json:"weight"`
}

type config struct {
	SOCKS5ListenAddr string            `json:"socks5_listen_addr"`
	Interfaces       []interfaceConfig `json:"interfaces"`
}

func loadConfig(filename string) (config, error) {
	fb, err := ioutil.ReadFile(filename)
	if err != nil {
		return config{}, err
	}
	var currentConfig config
	if err := json.Unmarshal(fb, &currentConfig); err != nil {
		return config{}, err
	}
	return currentConfig, err
}

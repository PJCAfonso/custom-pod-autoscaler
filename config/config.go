/*
Copyright 2019 The Custom Pod Autoscaler Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configEnvName   = "CONFIG_PATH"
	evaluateEnvName = "EVALUATE"
	metricEnvName   = "METRIC"
	intervalEnvName = "INTERVAL"
	selectorEnvName = "SELECTOR"
)

const (
	defaultConfig   = "/config.yaml"
	defaultEvaluate = ">&2 echo 'ERROR: No evaluate command set' && exit 1"
	defaultMetric   = ">&2 echo 'ERROR: No metric command set' && exit 1"
	defaultInterval = 15000
	defaultSelector = ""
)

// Config is the configuration options for the CPA
type Config struct {
	Evaluate string `yaml:"evaluate"`
	Metric   string `yaml:"metric"`
	Interval int    `yaml:"interval"`
	Selector string `yaml:"selector"`
}

// LoadConfig loads in the default configuration, then overrides it from the config file,
// then any env vars set.
func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(getEnv(configEnvName, defaultConfig))
	if err != nil {
		return nil, err
	}
	config := newDefaultConfig()
	err = loadFromYAML(data, config)
	if err != nil {
		return nil, err
	}
	err = loadFromEnv(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadFromYAML(data []byte, config *Config) error {
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

func loadFromEnv(config *Config) error {
	// Get string env vars
	config.Selector = getEnv(selectorEnvName, config.Selector)
	config.Evaluate = getEnv(evaluateEnvName, config.Evaluate)
	config.Metric = getEnv(metricEnvName, config.Metric)
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func newDefaultConfig() *Config {
	return &Config{
		Interval: defaultInterval,
		Metric:   defaultMetric,
		Evaluate: defaultEvaluate,
		Selector: defaultSelector,
	}
}

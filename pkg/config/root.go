package config

import "gopkg.in/yaml.v2"

type Root struct {
	Runtime string `yaml: "runtime"`
	Nodes   []Node `yaml:"nodes"`
}

func ParseConfigRoot(configRoot string) (*Root, error) {
	root := &Root{}
	err := yaml.Unmarshal([]byte(configRoot), root)
	if err != nil {
		return nil, err
	}
	return root, nil
}

package config

type Node struct {
	Cluster string `yaml: "cluster"`
	Id      string `yaml: "id"`
}

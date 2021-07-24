package main

import (
	"fmt"
	"os"

	"github.com/gotway/gotway/pkg/env"
)

type config struct {
	kubeConfig string
	namespace  string
	numWorkers int
	ha         bool
	nodeId     string
	env        string
	logLevel   string
}

func (c config) String() string {
	return fmt.Sprintf(
		"config{kubeConfig='%s'namespace='%s'numWorkers='%d'ha='%v'nodeId='%s'env='%s'loglevel='%s'}",
		c.kubeConfig,
		c.namespace,
		c.numWorkers,
		c.ha,
		c.nodeId,
		c.env,
		c.logLevel,
	)
}

func getConfig() (config, error) {
	ha := env.GetBool("HA", false)

	var nodeId string
	if ha {
		nodeId = env.Get("NODE_ID", "")
		if nodeId == "" {
			hostname, err := os.Hostname()
			if err != nil {
				return config{}, fmt.Errorf("error getting node id %v", err)
			}
			nodeId = hostname
		}
	}

	return config{
		kubeConfig: env.Get("KUBECONFIG", ""),
		namespace:  env.Get("NAMESPACE", "default"),
		numWorkers: env.GetInt("NUM_WORKERS", 4),
		ha:         env.GetBool("HA", false),
		nodeId:     nodeId,
		env:        env.Get("ENV", "local"),
		logLevel:   env.Get("LOG_LEVEL", "debug"),
	}, nil
}

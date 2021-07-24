package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gotway/gotway/pkg/env"
)

type haConfig struct {
	enabled       bool
	nodeId        string
	leaseLockName string
	leaseDuration time.Duration
	renewDeadline time.Duration
	retryPeriod   time.Duration
}

func (c haConfig) String() string {
	return fmt.Sprintf(
		"haConfig{enabled='%v'nodeId='%s'leaseLockName='%s'leaseDuration='%v'renewDeadline='%v'retryPeriod='%v'}",
		c.enabled,
		c.nodeId,
		c.leaseLockName,
		c.leaseDuration,
		c.renewDeadline,
		c.retryPeriod,
	)
}

type config struct {
	kubeConfig string
	namespace  string
	numWorkers int
	ha         haConfig
	env        string
	logLevel   string
}

func (c config) String() string {
	return fmt.Sprintf(
		"config{kubeConfig='%s'namespace='%s'numWorkers='%d'ha='%v'env='%s'loglevel='%s'}",
		c.kubeConfig,
		c.namespace,
		c.numWorkers,
		c.ha,
		c.env,
		c.logLevel,
	)
}

func getConfig() (config, error) {
	ha := env.GetBool("HA_ENABLED", false)

	var nodeId string
	if ha {
		nodeId = env.Get("HA_NODE_ID", "")
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
		ha: haConfig{
			enabled:       ha,
			nodeId:        nodeId,
			leaseLockName: env.Get("HA_LEASE_LOCK_NAME", "echoperator"),
			leaseDuration: env.GetDuration("HA_LEASE_DURATION_SECONDS", 15) * time.Second,
			renewDeadline: env.GetDuration("HA_RENEW_DEADLINE_SECONDS", 10) * time.Second,
			retryPeriod:   env.GetDuration("HA_RETRY_PERIOD_SECONDS", 2) * time.Second,
		},
		env:      env.Get("ENV", "local"),
		logLevel: env.Get("LOG_LEVEL", "debug"),
	}, nil
}

package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gotway/gotway/pkg/env"
)

type HAConfig struct {
	Enabled       bool
	NodeId        string
	LeaseLockName string
	LeaseDuration time.Duration
	RenewDeadline time.Duration
	RetryPeriod   time.Duration
}

func (c HAConfig) String() string {
	return fmt.Sprintf(
		"HAConfig{Enabled='%v'NodeId='%s'LeaseLockName='%s'LeaseDuration='%v'RenewDeadline='%v'RetryPeriod='%v'}",
		c.Enabled,
		c.NodeId,
		c.LeaseLockName,
		c.LeaseDuration,
		c.RenewDeadline,
		c.RetryPeriod,
	)
}

type Config struct {
	KubeConfig string
	Namespace  string
	NumWorkers int
	HA         HAConfig
	Env        string
	LogLevel   string
}

func (c Config) String() string {
	return fmt.Sprintf(
		"Config{KubeConfig='%s'Namespace='%s'NumWorkers='%d'HA='%v'Env='%s'LogLevel='%s'}",
		c.KubeConfig,
		c.Namespace,
		c.NumWorkers,
		c.HA,
		c.Env,
		c.LogLevel,
	)
}

func GetConfig() (Config, error) {
	ha := env.GetBool("HA_ENABLED", false)

	var nodeId string
	if ha {
		nodeId = env.Get("HA_NODE_ID", "")
		if nodeId == "" {
			hostname, err := os.Hostname()
			if err != nil {
				return Config{}, fmt.Errorf("error getting node id %v", err)
			}
			nodeId = hostname
		}
	}

	return Config{
		KubeConfig: env.Get("KUBECONFIG", ""),
		Namespace:  env.Get("NAMESPACE", "default"),
		NumWorkers: env.GetInt("NUM_WORKERS", 4),
		HA: HAConfig{
			Enabled:       ha,
			NodeId:        nodeId,
			LeaseLockName: env.Get("HA_LEASE_LOCK_NAME", "echoperator"),
			LeaseDuration: env.GetDuration("HA_LEASE_DURATION_SECONDS", 15) * time.Second,
			RenewDeadline: env.GetDuration("HA_RENEW_DEADLINE_SECONDS", 10) * time.Second,
			RetryPeriod:   env.GetDuration("HA_RETRY_PERIOD_SECONDS", 2) * time.Second,
		},
		Env:      env.Get("ENV", "local"),
		LogLevel: env.Get("LOG_LEVEL", "debug"),
	}, nil
}

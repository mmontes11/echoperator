package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/metrics"

	"github.com/mmontes11/echoperator/internal/config"
	"github.com/mmontes11/echoperator/internal/runner"
	"github.com/mmontes11/echoperator/pkg/controller"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error getting config %v", err))
	}
	logger := getLogger(config)
	logger.Debugf("config %v", config)

	var restConfig *rest.Config
	var errKubeConfig error
	if config.KubeConfig != "" {
		restConfig, errKubeConfig = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	} else {
		restConfig, errKubeConfig = rest.InClusterConfig()
	}
	if errKubeConfig != nil {
		logger.Fatal("error getting kubernetes config ", err)
	}

	kubeClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		logger.Fatal("error getting kubernetes client ", err)
	}
	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		logger.Fatal("error creating api extensions client ", err)
	}
	echov1alpha1ClientSet, err := echov1alpha1clientset.NewForConfig(restConfig)
	if err != nil {
		logger.Fatal("error creating echo client ", err)
	}

	ctrl := controller.New(
		kubeClientSet,
		apiextensionsClientSet,
		echov1alpha1ClientSet,
		config.Namespace,
		logger.WithField("type", "controller"),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	defer cancel()

	if config.Metrics.Enabled {
		m := metrics.New(
			metrics.Options{
				Path: config.Metrics.Path,
				Port: config.Metrics.Port,
			},
			logger.WithField("type", "metrics"),
		)
		go m.Start()
		defer m.Stop()
	}

	r := runner.NewRunner(
		ctrl,
		kubeClientSet,
		config,
		logger.WithField("type", "runner"),
	)
	r.Start(ctx)
}

func getLogger(config config.Config) log.Logger {
	logger := log.NewLogger(log.Fields{
		"service": "echoperator",
	}, config.Env, config.LogLevel, os.Stdout)
	if config.HA.Enabled {
		return logger.WithField("node", config.HA.NodeId)
	}
	return logger
}

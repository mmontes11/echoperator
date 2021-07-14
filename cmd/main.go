package main

import (
	"context"
	"os"

	gs "github.com/gotway/gotway/pkg/gracefulshutdown"
	"github.com/gotway/gotway/pkg/log"
	crd "github.com/mmontes11/echoperator/pkg/kubernetes/echoperator"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	logger := log.NewLogger(log.Fields{
		"service": "echoperator",
	}, "local", "debug", os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var config *rest.Config
	var err error
	if kubeconfig, ok := os.LookupEnv("KUBECONFIG"); ok {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		logger.Fatal("error getting kubernetes config ", err)
	}

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		logger.Fatal("error creating api extensions client ", err)
	}

	logger.Info("creating curstom resource definition")
	_, err = crd.CreateCustomResourceDefinition(ctx, apiextensionsClientSet)
	if err != nil {
		logger.Fatal("error creating custom resource definition ", err)
	}

	gs.GracefulShutdown(logger, cancel)
}

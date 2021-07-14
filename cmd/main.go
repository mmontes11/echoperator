package main

import (
	"context"
	"os"

	gs "github.com/gotway/gotway/pkg/gracefulshutdown"
	"github.com/gotway/gotway/pkg/log"

	"github.com/mmontes11/echoperator/pkg/controller"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
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

	kubeClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("error getting kubernetes client ", err)
	}

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		logger.Fatal("error creating api extensions client ", err)
	}

	echov1alpha1clientset, err := echov1alpha1clientset.NewForConfig(config)
	if err != nil {
		logger.Fatal("error creating echo client ", err)
	}

	ctrl := controller.New(
		kubeClientSet,
		apiextensionsClientSet,
		echov1alpha1clientset,
		"default",
		"default",
	)

	_, err = ctrl.RegisterCustomResourceDefinition(ctx)
	if err != nil {
		logger.Fatal("error registering custom resource definition ", err)
	}
	logger.Info("custom resource definition registered")

	gs.GracefulShutdown(logger, cancel)
}

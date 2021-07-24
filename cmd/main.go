package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/pkg/log"

	"github.com/mmontes11/echoperator/pkg/controller"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := getConfig()
	if err != nil {
		panic(fmt.Errorf("error getting config %v", err))
	}
	ctx, _ := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	logger := log.NewLogger(log.Fields{
		"service": "echoperator",
	}, config.env, config.logLevel, os.Stdout)

	logger.Debugf("config %v", config)

	var restConfig *rest.Config
	var errKubeConfig error
	if config.kubeConfig != "" {
		restConfig, errKubeConfig = clientcmd.BuildConfigFromFlags("", config.kubeConfig)
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
		config.namespace,
		logger.WithField("type", "controller"),
	)

	if err := ctrl.RegisterCustomResourceDefinitions(ctx); err != nil {
		logger.Fatal("error registering CRDs ", err)
	}

	if err := ctrl.Run(ctx, config.numWorkers); err != nil {
		logger.Fatal("error running controller ", err)
	}
}

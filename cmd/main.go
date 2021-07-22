package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/pkg/env"
	"github.com/gotway/gotway/pkg/log"

	"github.com/mmontes11/echoperator/pkg/controller"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig = env.Get("KUBECONFIG", "")
	Namespace  = env.Get("NAMESPACE", "default")
	NumWorkers = env.GetInt("NUM_WORKERS", 4)
	Env        = env.Get("ENV", "local")
	LogLevel   = env.Get("LOG_LEVEL", "debug")
)

func main() {
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
	}, Env, LogLevel, os.Stdout)

	var config *rest.Config
	var err error
	if KubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", KubeConfig)
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
	echov1alpha1ClientSet, err := echov1alpha1clientset.NewForConfig(config)
	if err != nil {
		logger.Fatal("error creating echo client ", err)
	}

	ctrl := controller.New(
		kubeClientSet,
		apiextensionsClientSet,
		echov1alpha1ClientSet,
		Namespace,
		logger.WithField("type", "controller"),
	)

	if err := ctrl.RegisterCustomResourceDefinitions(ctx); err != nil {
		logger.Fatal("error registering CRDs ", err)
	}

	if err := ctrl.Run(ctx, NumWorkers); err != nil {
		logger.Fatal("error running controller ", err)
	}
}

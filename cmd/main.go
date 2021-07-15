package main

import (
	"context"
	"os"
	"time"

	"github.com/gotway/gotway/pkg/env"
	"github.com/gotway/gotway/pkg/log"

	echoperatorctx "github.com/mmontes11/echoperator/pkg/context"
	"github.com/mmontes11/echoperator/pkg/controller"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"
	echoinformers "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/informers/externalversions"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig          = env.Get("KUBECONFIG", "")
	KubernetesNamespace = env.Get("KUBERNETES_NAMESPACE", "default")
	CRDNamespace        = env.Get("CRD_NAMESPACE", "default")
	NumWorkers          = env.GetInt("NUM_WORKERS", 4)
)

func main() {
	logger := log.NewLogger(log.Fields{
		"service": "echoperator",
	}, "local", "debug", os.Stdout)
	ctx := echoperatorctx.WithGracefulShutdown(context.Background(), logger)

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

	echoInformerFactory := echoinformers.NewSharedInformerFactory(echov1alpha1ClientSet, time.Second*30)

	ctrl := controller.New(
		kubeClientSet,
		apiextensionsClientSet,
		echov1alpha1ClientSet,
		echoInformerFactory.Mmontes().V1alpha1().Echos(),
		KubernetesNamespace,
		CRDNamespace,
		logger.WithField("type", "controller"),
	)

	_, err = ctrl.RegisterCustomResourceDefinition(ctx)
	if err != nil {
		logger.Fatal("error registering custom resource definition ", err)
	}

	if err := ctrl.Run(ctx, NumWorkers); err != nil {
		logger.Fatal("error running controller ", err)
	}
}

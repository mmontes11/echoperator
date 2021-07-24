package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/echoperator/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var sigs = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGKILL,
	syscall.SIGHUP,
	syscall.SIGQUIT,
}

func run(
	ctrl *controller.Controller,
	clientset *kubernetes.Clientset,
	config config,
	logger log.Logger,
) {
	if config.ha {
		runHA(ctrl, clientset, config, logger)
	} else {
		runStandalone(ctrl, config, logger)
	}
}

func runSingleNode(
	ctx context.Context,
	ctrl *controller.Controller,
	config config,
	logger log.Logger,
) {
	if err := ctrl.RegisterCustomResourceDefinitions(ctx); err != nil {
		logger.Fatal("error registering CRDs ", err)
	}
	if err := ctrl.Run(ctx, config.numWorkers); err != nil {
		logger.Fatal("error running controller ", err)
	}
}

func runStandalone(
	ctrl *controller.Controller,
	config config,
	logger log.Logger,
) {
	ctx, _ := signal.NotifyContext(context.Background(), sigs...)
	runSingleNode(ctx, ctrl, config, logger)
}

func runHA(
	ctrl *controller.Controller,
	clientset *kubernetes.Clientset,
	config config,
	logger log.Logger,
) {
	ctx, _ := signal.NotifyContext(context.Background(), sigs...)

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "echoperator",
			Namespace: config.namespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: config.nodeId,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				logger.Info("start leading")
				runSingleNode(ctx, ctrl, config, logger)
			},
			OnStoppedLeading: func() {
				logger.Info("stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == config.nodeId {
					logger.Info("obtained leadership")
					return
				}
				logger.Infof("leader elected: '%s'", identity)
			},
		},
	})
}

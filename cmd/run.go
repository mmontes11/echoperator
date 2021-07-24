package main

import (
	"context"
	"time"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/echoperator/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type runner struct {
	ctrl      *controller.Controller
	clientset *kubernetes.Clientset
	config    config
	logger    log.Logger
}

func (r *runner) start(ctx context.Context) {
	if r.config.ha {
		r.logger.Info("starting HA controller")
		r.runHA(ctx)
	} else {
		r.logger.Info("starting standalone controller")
		r.runSingleNode(ctx)
	}
}

func (r *runner) runSingleNode(ctx context.Context) {
	if err := r.ctrl.Run(ctx, r.config.numWorkers); err != nil {
		r.logger.Fatal("error running controller ", err)
	}
}

func (r *runner) runHA(ctx context.Context) {
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "echoperator",
			Namespace: r.config.namespace,
		},
		Client: r.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: r.config.nodeId,
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
				r.logger.Info("start leading")
				r.runSingleNode(ctx)
			},
			OnStoppedLeading: func() {
				r.logger.Info("stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == r.config.nodeId {
					r.logger.Info("obtained leadership")
					return
				}
				r.logger.Infof("leader elected: '%s'", identity)
			},
		},
	})
}

func newRunner(
	ctrl *controller.Controller,
	clientset *kubernetes.Clientset,
	config config,
	logger log.Logger,
) *runner {
	return &runner{
		ctrl:      ctrl,
		clientset: clientset,
		config:    config,
		logger:    logger,
	}
}

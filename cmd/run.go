package main

import (
	"context"

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
	if r.config.ha.enabled {
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
	if r.config.ha == (haConfig{}) || !r.config.ha.enabled {
		r.logger.Fatal("HA config not set or not enabled")
	}

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      r.config.ha.leaseLockName,
			Namespace: r.config.namespace,
		},
		Client: r.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: r.config.ha.nodeId,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   r.config.ha.leaseDuration,
		RenewDeadline:   r.config.ha.renewDeadline,
		RetryPeriod:     r.config.ha.retryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				r.logger.Info("start leading")
				r.runSingleNode(ctx)
			},
			OnStoppedLeading: func() {
				r.logger.Info("stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == r.config.ha.nodeId {
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

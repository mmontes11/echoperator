package runner

import (
	"context"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/echoperator/internal/config"
	"github.com/mmontes11/echoperator/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type Runner struct {
	ctrl      *controller.Controller
	clientset *kubernetes.Clientset
	config    config.Config
	logger    log.Logger
}

func (r *Runner) Start(ctx context.Context) {
	if r.config.HA.Enabled {
		r.logger.Info("starting HA controller")
		r.runHA(ctx)
	} else {
		r.logger.Info("starting standalone controller")
		r.runSingleNode(ctx)
	}
}

func (r *Runner) runSingleNode(ctx context.Context) {
	if err := r.ctrl.Run(ctx, r.config.NumWorkers); err != nil {
		r.logger.Fatal("error running controller ", err)
	}
}

func (r *Runner) runHA(ctx context.Context) {
	if r.config.HA == (config.HA{}) || !r.config.HA.Enabled {
		r.logger.Fatal("HA config not set or not enabled")
	}

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      r.config.HA.LeaseLockName,
			Namespace: r.config.Namespace,
		},
		Client: r.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: r.config.HA.NodeId,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   r.config.HA.LeaseDuration,
		RenewDeadline:   r.config.HA.RenewDeadline,
		RetryPeriod:     r.config.HA.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				r.logger.Info("start leading")
				r.runSingleNode(ctx)
			},
			OnStoppedLeading: func() {
				r.logger.Info("stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == r.config.HA.NodeId {
					r.logger.Info("obtained leadership")
					return
				}
				r.logger.Infof("leader elected: '%s'", identity)
			},
		},
	})
}

func NewRunner(
	ctrl *controller.Controller,
	clientset *kubernetes.Clientset,
	config config.Config,
	logger log.Logger,
) *Runner {
	return &Runner{
		ctrl:      ctrl,
		clientset: clientset,
		config:    config,
		logger:    logger,
	}
}

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gotway/gotway/pkg/log"

	echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"
	echoinformers "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/informers/externalversions"

	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeClientSet kubernetes.Interface
	extClientSet  extclientset.Interface

	echoInformer          cache.SharedIndexInformer
	jobInformer           cache.SharedIndexInformer
	scheduledEchoInformer cache.SharedIndexInformer
	cronjobInformer       cache.SharedIndexInformer

	queue workqueue.RateLimitingInterface

	namespace string

	logger log.Logger
}

func (c *Controller) Run(ctx context.Context, numWorkers int) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("starting controller")

	c.logger.Info("starting informers")
	for _, i := range []cache.SharedIndexInformer{
		c.echoInformer,
		c.scheduledEchoInformer,
		c.jobInformer,
		c.cronjobInformer,
	} {
		go i.Run(ctx.Done())
	}

	c.logger.Info("waiting for informer caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), []cache.InformerSynced{
		c.echoInformer.HasSynced,
		c.scheduledEchoInformer.HasSynced,
		c.jobInformer.HasSynced,
		c.cronjobInformer.HasSynced,
	}...) {
		return fmt.Errorf("failed to wait for informers caches to sync")
	}

	c.logger.Infof("starting %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go wait.Until(func() { c.runWorker(ctx) }, time.Second, ctx.Done())
	}
	c.logger.Info("controller ready")

	<-ctx.Done()
	c.logger.Info("stopping controller")

	return nil
}

func (c *Controller) AddEcho(obj interface{}) {
	c.logger.Debug("adding echo")
	echo, ok := obj.(*echov1alpha1.Echo)
	if !ok {
		c.logger.Errorf("unexpected object %v", obj)
		return
	}

	c.queue.Add(event{
		eventType: addEcho,
		resource:  echo.DeepCopy(),
	})
}

func (c *Controller) AddScheduledEcho(obj interface{}) {
	c.logger.Debug("adding scheduled echo")
	scheduledEcho, ok := obj.(*echov1alpha1.ScheduledEcho)
	if !ok {
		c.logger.Errorf("unexpected object %v", obj)
		return
	}

	c.queue.Add(event{
		eventType: addScheduledEcho,
		resource:  scheduledEcho.DeepCopy(),
	})
}

func New(
	kubeClientSet kubernetes.Interface,
	extClientSet extclientset.Interface,
	echoClientSet echov1alpha1clientset.Interface,
	namespace string,
	logger log.Logger,
) *Controller {

	echoInformerFactory := echoinformers.NewSharedInformerFactory(
		echoClientSet,
		10*time.Second,
	)
	echoInformer := echoInformerFactory.Mmontes().V1alpha1().Echos().Informer()
	scheduledechoInformer := echoInformerFactory.Mmontes().V1alpha1().ScheduledEchos().Informer()

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, 10*time.Second)
	jobInformer := kubeInformerFactory.Batch().V1().Jobs().Informer()
	cronjobInformer := kubeInformerFactory.Batch().V1beta1().CronJobs().Informer()

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	ctrl := &Controller{
		kubeClientSet: kubeClientSet,
		extClientSet:  extClientSet,

		echoInformer:          echoInformer,
		jobInformer:           jobInformer,
		scheduledEchoInformer: scheduledechoInformer,
		cronjobInformer:       cronjobInformer,

		queue: queue,

		namespace: namespace,

		logger: logger,
	}

	echoInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.AddEcho,
	})
	scheduledechoInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.AddScheduledEcho,
	})

	return ctrl
}

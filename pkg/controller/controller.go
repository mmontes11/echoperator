package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gotway/gotway/pkg/log"

	echo "github.com/mmontes11/echoperator/pkg/echo"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"
	echoinformer "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/informers/externalversions/echoperator/v1alpha1"
	echolister "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/listers/echoperator/v1alpha1"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeClientSet kubernetes.Interface
	extClientSet  extclientset.Interface

	echoClientSet echov1alpha1clientset.Interface
	echoLister    echolister.EchoLister
	echoInformer  cache.SharedIndexInformer

	kubeNamespace string
	crdNamespace  string

	queue workqueue.RateLimitingInterface

	logger log.Logger
}

func (c *Controller) RegisterCustomResourceDefinition(
	ctx context.Context,
) (extv1.CustomResourceDefinition, error) {

	_, err := c.extClientSet.ApiextensionsV1().
		CustomResourceDefinitions().
		Create(ctx, &echo.CRD, metav1.CreateOptions{})

	if err != nil && !apierrors.IsAlreadyExists(err) {
		return extv1.CustomResourceDefinition{}, err
	}

	err = wait.Poll(5*time.Second, 1*time.Minute, func() (bool, error) {
		crd, err := c.extClientSet.ApiextensionsV1().
			CustomResourceDefinitions().
			Get(ctx, echo.CRDName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			if cond.Type == extv1.Established &&
				cond.Status == extv1.ConditionTrue {
				return true, nil
			}
		}
		return false, err
	})

	if err != nil {
		deleteErr := c.extClientSet.ApiextensionsV1().
			CustomResourceDefinitions().
			Delete(ctx, echo.CRDName, metav1.DeleteOptions{})
		if deleteErr != nil {
			return extv1.CustomResourceDefinition{}, errors.NewAggregate([]error{err, deleteErr})
		}
	}
	return echo.CRD, nil
}

func (c *Controller) Run(ctx context.Context, numWorkers int) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("starting controller")

	c.logger.Info("starting informer")
	go c.echoInformer.Run(ctx.Done())

	c.logger.Info("waiting for informer caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), c.echoInformer.HasSynced) {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	c.logger.Infof("starting %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go wait.Until(c.runWorker, time.Second, ctx.Done())
	}
	c.logger.Info("controller ready")

	<-ctx.Done()
	c.logger.Info("stopping controller")

	return nil
}

func (c *Controller) Add(obj interface{}) {
	c.logger.Info("adding")
}

func (c *Controller) Update(oldObj interface{}, newObj interface{}) {
	c.logger.Info("updating")
}

func (c *Controller) Delete(obj interface{}) {
	c.logger.Info("deleting")
}

func New(
	kubeClientSet kubernetes.Interface,
	extClientSet extclientset.Interface,
	echoClientSet echov1alpha1clientset.Interface,
	echoInformer echoinformer.EchoInformer,
	kubeNamespace, echoNamespace string,
	logger log.Logger,
) *Controller {

	informer := echoInformer.Informer()

	ctrl := &Controller{
		kubeClientSet: kubeClientSet,
		extClientSet:  extClientSet,

		echoClientSet: echoClientSet,
		echoLister:    echoInformer.Lister(),
		echoInformer:  informer,

		kubeNamespace: kubeNamespace,
		crdNamespace:  echoNamespace,

		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),

		logger: logger,
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctrl.Add,
		UpdateFunc: ctrl.Update,
		DeleteFunc: ctrl.Delete,
	})

	return ctrl
}

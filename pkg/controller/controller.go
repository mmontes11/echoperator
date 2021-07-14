package controller

import (
	"context"
	"time"

	echo "github.com/mmontes11/echoperator/pkg/echo"
	echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	echov1alpha1clientset "github.com/mmontes11/echoperator/pkg/echo/v1alpha1/apis/clientset/versioned"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeClientSet kubernetes.Interface
	extClientSet  extclientset.Interface
	echoClientSet echov1alpha1clientset.Interface

	kubeNamespace string
	crdNamespace  string

	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
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

func New(
	kubeClientSet kubernetes.Interface,
	extClientSet extclientset.Interface,
	echoClientSet echov1alpha1clientset.Interface,
	kubeNamespace, echoNamespace string,
) *Controller {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	lw := cache.NewListWatchFromClient(
		echoClientSet.MmontesV1alpha1().RESTClient(),
		echo.Plural,
		echoNamespace,
		fields.Everything(),
	)
	informer := cache.NewSharedIndexInformer(lw, &echov1alpha1.Echo{}, 0, cache.Indexers{})

	return &Controller{
		kubeClientSet: kubeClientSet,
		extClientSet:  extClientSet,
		echoClientSet: echoClientSet,
		kubeNamespace: kubeNamespace,
		crdNamespace:  echoNamespace,
		queue:         queue,
		informer:      informer,
	}
}

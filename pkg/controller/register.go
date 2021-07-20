package controller

import (
	"context"
	"fmt"
	"time"

	echo "github.com/mmontes11/echoperator/pkg/echo"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (c *Controller) RegisterCustomResourceDefinitions(ctx context.Context) error {
	crds := []extv1.CustomResourceDefinition{echo.EchoCRD, echo.ScheduledEchoCRD}
	for _, crd := range crds {
		if err := c.registerCustomResourceDefinition(ctx, &crd); err != nil {
			return fmt.Errorf(
				"error registering custom resource definition '%s': %v",
				crd.ObjectMeta.Name,
				err,
			)
		}
	}
	return nil
}

func (c *Controller) registerCustomResourceDefinition(
	ctx context.Context,
	crd *extv1.CustomResourceDefinition,
) error {

	_, err := c.extClientSet.ApiextensionsV1().
		CustomResourceDefinitions().
		Create(ctx, crd, metav1.CreateOptions{})

	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}

	err = wait.Poll(5*time.Second, 1*time.Minute, func() (bool, error) {
		crd, err := c.extClientSet.ApiextensionsV1().
			CustomResourceDefinitions().
			Get(ctx, crd.ObjectMeta.Name, metav1.GetOptions{})
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
			Delete(ctx, crd.ObjectMeta.Name, metav1.DeleteOptions{})
		if deleteErr != nil {
			return errors.NewAggregate([]error{err, deleteErr})
		}
	}
	return nil
}

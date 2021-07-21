package controller

import (
	"context"
	"fmt"

	echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const maxRetries = 3

func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextItem(ctx) {
	}
}

func (c *Controller) processNextItem(ctx context.Context) bool {
	obj, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(obj)

	err := c.processEvent(ctx, obj)
	if err == nil {
		c.logger.Debug("processed item")
		c.queue.Forget(obj)
	} else if c.queue.NumRequeues(obj) < maxRetries {
		c.logger.Errorf("error processing event: %v, retrying", err)
		c.queue.AddRateLimited(obj)
	} else {
		c.logger.Errorf("error processing event: %v, max retries reached", err)
		c.queue.Forget(obj)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processEvent(ctx context.Context, obj interface{}) error {
	event, ok := obj.(event)
	if !ok {
		c.logger.Error("unexpected event ", obj)
		return nil
	}
	switch event.eventType {
	case addEcho:
		return c.addEcho(ctx, event.resource.(*echov1alpha1.Echo))
	case addScheduledEcho:
		return c.addScheduledEcho(ctx, event.resource.(*echov1alpha1.ScheduledEcho))
	}
	return nil
}

func (c *Controller) addEcho(ctx context.Context, echo *echov1alpha1.Echo) error {
	exists, err := c.jobAlreadyExists(echo)
	if err != nil {
		return err
	}
	if exists {
		c.logger.Debug("echo job already exists, skipping")
		return nil
	}
	job := createJob(echo, c.namespace)
	_, err = c.kubeClientSet.BatchV1().
		Jobs(c.namespace).
		Create(ctx, job, metav1.CreateOptions{})
	return err
}

func (c *Controller) addScheduledEcho(ctx context.Context, scheduledEcho *echov1alpha1.ScheduledEcho) error {
	exists, err := c.cronJobAlreadyExists(scheduledEcho)
	if err != nil {
		return err
	}
	if exists {
		c.logger.Debug("echo cronjob already exists, skipping")
		return nil
	}
	cronjob := createCronJob(scheduledEcho, c.namespace)
	_, err = c.kubeClientSet.BatchV1beta1().
		CronJobs(c.namespace).
		Create(ctx, cronjob, metav1.CreateOptions{})
	return err
}

func (c *Controller) jobAlreadyExists(echo *echov1alpha1.Echo) (bool, error) {
	for _, obj := range c.jobInformer.GetIndexer().List() {
		job, ok := (obj).(*batchv1.Job)
		if !ok {
			return false, fmt.Errorf("unexpected object %v", obj)
		}
		for _, owner := range job.ObjectMeta.OwnerReferences {
			if owner.UID == echo.UID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (c *Controller) cronJobAlreadyExists(echo *echov1alpha1.ScheduledEcho) (bool, error) {
	for _, obj := range c.cronjobInformer.GetIndexer().List() {
		cronjob, ok := (obj).(*batchv1beta1.CronJob)
		if !ok {
			return false, fmt.Errorf("unexpected object %v", obj)
		}
		for _, owner := range cronjob.ObjectMeta.OwnerReferences {
			if owner.UID == echo.UID {
				return true, nil
			}
		}
	}
	return false, nil
}

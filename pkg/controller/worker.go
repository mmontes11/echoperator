package controller

import (
	"context"

	echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

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

	event, ok := obj.(event)
	if !ok {
		c.logger.Errorf("unexpected event %v", event)
		return true
	}

	err := c.processEvent(ctx, event)
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

func (c *Controller) processEvent(ctx context.Context, event event) error {
	switch event.eventType {
	case add:
		return c.addEcho(ctx, event.newEcho)
	case delete:
	case update:
	}
	return nil
}

func (c *Controller) addEcho(ctx context.Context, echo *echov1alpha1.Echo) error {
	job := createJob(echo, c.kubeNamespace)
	_, err := c.kubeClientSet.BatchV1().
		Jobs(c.kubeNamespace).
		Create(ctx, job, metav1.CreateOptions{})
	return err
}

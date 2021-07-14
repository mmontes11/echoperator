package controller

import (
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	kubeClientSet kubernetes.Interface
}

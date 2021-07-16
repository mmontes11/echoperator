package controller

import (
	echo "github.com/mmontes11/echoperator/pkg/echo"
	echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	"github.com/mmontes11/echoperator/pkg/echo/version"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createJob(newEcho *echov1alpha1.Echo, namespace string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: newEcho.ObjectMeta.Name + "-echo-",
			Namespace:    namespace,
			Labels:       make(map[string]string),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: version.V1alpha1,
					Kind:       echo.Kind,
					Name:       newEcho.ObjectMeta.Name,
					UID:        newEcho.UID,
				},
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: newEcho.Name + "-",
					Namespace:    namespace,
					Labels:       make(map[string]string),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            newEcho.Name,
							Image:           "busybox:1.33.1",
							Command:         []string{"echo", newEcho.Spec.Message},
							ImagePullPolicy: "IfNotPresent",
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}

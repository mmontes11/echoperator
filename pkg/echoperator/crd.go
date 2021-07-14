package v1alpha1

import (
	"context"
	"time"

	version "github.com/mmontes11/echoperator/pkg/echoperator/version"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateCustomResourceDefinition(
	ctx context.Context,
	clientSet clientset.Interface,
) (*extv1.CustomResourceDefinition, error) {
	crd := &extv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: CRDName,
		},
		Spec: extv1.CustomResourceDefinitionSpec{
			Group: GroupName,
			Names: extv1.CustomResourceDefinitionNames{
				Kind:       Kind,
				Plural:     Plural,
				Singular:   Singular,
				ShortNames: ShortNames,
			},
			Scope: extv1.NamespaceScoped,
			Versions: []extv1.CustomResourceDefinitionVersion{
				{
					Name:    version.V1alpha1,
					Served:  true,
					Storage: true,
					Schema: &extv1.CustomResourceValidation{
						OpenAPIV3Schema: &extv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]extv1.JSONSchemaProps{
								"message": {Type: "string"},
							},
							Required: []string{"message"},
						},
					},
				},
			},
		},
	}
	_, err := clientSet.ApiextensionsV1().
		CustomResourceDefinitions().
		Create(ctx, crd, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, err
	}

	err = wait.Poll(5*time.Second, 1*time.Minute, func() (bool, error) {
		crd, err := clientSet.ApiextensionsV1().
			CustomResourceDefinitions().
			Get(ctx, CRDName, metav1.GetOptions{})
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
		deleteErr := clientSet.ApiextensionsV1().
			CustomResourceDefinitions().
			Delete(ctx, CRDName, metav1.DeleteOptions{})
		if deleteErr != nil {
			return nil, errors.NewAggregate([]error{err, deleteErr})
		}
	}

	return crd, nil
}

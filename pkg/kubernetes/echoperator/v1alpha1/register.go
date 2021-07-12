package v1alpha1

import (
	crd "github.com/mmontes11/echoperator/pkg/kubernetes/echoperator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Kind         string = "Echo"
	GroupVersion string = "v1alpha1"
	Plural       string = "echo"
	Singular     string = "echos"
	CRDName      string = Plural + "." + crd.GroupName
	ShortName    string = "ec"
)

var (
	SchemeGroupVersion = schema.GroupVersion{
		Group:   crd.GroupName,
		Version: GroupVersion,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Echo{},
		&EchoList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

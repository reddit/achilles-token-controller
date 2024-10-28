// Package scheme contains the scheme for all GVKs interacted with by the controller.
package scheme

import (
	"fmt"

	"github.snooguts.net/reddit/achilles-token-controller/api/group.example.com/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	kscheme "k8s.io/client-go/kubernetes/scheme"
)

var AddToSchemes = runtime.SchemeBuilder{}

func init() {
	AddToSchemes.Register(kscheme.AddToScheme)  // native kubernetes schemes
	AddToSchemes.Register(v1alpha1.AddToScheme) // internal schemes
}

// NewScheme creates and populates a runtime.Scheme with the default k8s resources as well as Reddit's resources.
func NewScheme() (*runtime.Scheme, error) {
	s := runtime.NewScheme()

	// add all k8s native schemes
	if err := kscheme.AddToScheme(s); err != nil {
		return nil, fmt.Errorf("adding k8s resources to scheme: %w", err)
	}

	// add CRD schemes
	if err := AddToSchemes.AddToScheme(s); err != nil {
		return nil, fmt.Errorf("adding internal resources to scheme: %w", err)
	}

	return s, nil
}

func MustNewScheme() *runtime.Scheme {
	s, err := NewScheme()
	if err != nil {
		panic(err)
	}
	return s
}

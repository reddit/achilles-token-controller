package accesstoken

import (
	"github.com/reddit/achilles-sdk-api/api"
	"github.com/reddit/achilles-token-controller/api/group.example.com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var conditionTokenProvisioned = api.Condition{
	Type:    v1alpha1.TypeTokenProvisioned,
	Status:  corev1.ConditionTrue,
	Message: "Access token has been provisioned (see `status.tokenSecretRef`)",
}

var conditionStalePermissionsRemoved = api.Condition{
	Type:    v1alpha1.TypeStalePermissionsRemoved,
	Status:  corev1.ConditionTrue,
	Message: "Stale permissions have been removed",
}

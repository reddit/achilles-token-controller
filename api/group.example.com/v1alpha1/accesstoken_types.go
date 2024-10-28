package v1alpha1

import (
	"github.snooguts.net/reddit/achilles-sdk-api/api"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&AccessToken{}, &AccessTokenList{})
}

const (
	// TypeTokenProvisioned is a condition type that indicates the access token has been provisioned.
	TypeTokenProvisioned api.ConditionType = "TokenProvisioned"

	// TypeStalePermissionsRemoved is a condition type that indicates stale permissions have been removed.
	TypeStalePermissionsRemoved api.ConditionType = "StalePermissionsRemoved"
)

// AccessToken is the Schema for the AccessToken API
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
type AccessToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccessTokenSpec   `json:"spec,omitempty"`
	Status AccessTokenStatus `json:"status,omitempty"`
}

// AccessTokenList contains a list of AccessToken
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
type AccessTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AccessToken `json:"items"`
}

// AccessTokenSpec defines the desired state of AccessToken
type AccessTokenSpec struct {
	// NamespacedPermissions defines a list of namespaced scoped permissions. Optional
	NamespacedPermissions []NamespacedPermissions `json:"namespacedPermissions,omitempty"`

	// ClusterPermissions defines cluster scoped permissions. Optional
	ClusterPermissions *ClusterPermissions `json:"clusterPermissions,omitempty"`
}

type NamespacedPermissions struct {
	// Namespace the role applies to. Required
	Namespace string `json:"namespace"`

	// Rules for the role. Required
	Rules []rbacv1.PolicyRule `json:"rules"`
}

type ClusterPermissions struct {
	// Rules for the role. Required
	Rules []rbacv1.PolicyRule `json:"rules"`
}

// AccessTokenStatus defines the observed state of AccessToken
type AccessTokenStatus struct {
	api.ConditionedStatus `json:",inline"`

	// ResourceRefs is a list of all resources managed by this object.
	ResourceRefs []api.TypedObjectRef `json:"resourceRefs,omitempty"`

	// TokenSecretRef is a reference to the Secret containing the access token.
	TokenSecretRef *string `json:"tokenSecretRef,omitempty"`
}

func (c *AccessToken) GetConditions() []api.Condition {
	return c.Status.Conditions
}

func (c *AccessToken) SetConditions(cond ...api.Condition) {
	c.Status.SetConditions(cond...)
}

func (c *AccessToken) GetCondition(t api.ConditionType) api.Condition {
	return c.Status.GetCondition(t)
}

func (c *AccessToken) SetManagedResources(refs []api.TypedObjectRef) {
	c.Status.ResourceRefs = refs
}

func (c *AccessToken) GetManagedResources() []api.TypedObjectRef {
	return c.Status.ResourceRefs
}

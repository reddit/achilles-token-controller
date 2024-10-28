package accesstoken

import (
	"fmt"

	"github.snooguts.net/reddit/achilles-token-controller/api/group.example.com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type builder struct {
	accessToken *v1alpha1.AccessToken
}

func newBuilder(
	accessToken *v1alpha1.AccessToken,
) *builder {
	return &builder{
		accessToken: accessToken,
	}
}

func (b *builder) build() []client.Object {
	resources := []client.Object{
		b.serviceAccount(),
		b.secret(),
	}

	resources = append(resources, b.roleAndBindings()...)
	resources = append(resources, b.clusterRoleAndBinding()...)

	return resources
}

func (b *builder) serviceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.accessToken.GetName(),
			Namespace: b.accessToken.GetNamespace(),
		},
	}
}

func (b *builder) secret() *corev1.Secret {
	sa := b.serviceAccount()
	// https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#manually-create-a-long-lived-api-token-for-a-serviceaccount
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.accessToken.Name,
			Namespace: b.accessToken.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": sa.GetName(),
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}
}

func (b *builder) roleAndBindings() []client.Object {
	var objs []client.Object

	for _, namespacedRole := range b.accessToken.Spec.NamespacedPermissions {
		role := b.role(b.accessToken, namespacedRole.Namespace, namespacedRole.Rules)
		objs = append(objs, role)

		roleRef := rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     role.Name,
		}

		objs = append(objs, b.roleBinding(roleRef, namespacedRole.Namespace))
	}

	return objs
}

func (b *builder) role(accessToken *v1alpha1.AccessToken, ns string, rules []rbacv1.PolicyRule) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      accessToken.GetName(),
			Namespace: ns,
		},
		Rules: rules,
	}
}

func (b *builder) roleBinding(roleRef rbacv1.RoleRef, ns string) *rbacv1.RoleBinding {
	sa := b.serviceAccount()
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.accessToken.GetName(),
			Namespace: ns,
		},
		RoleRef: roleRef,
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      sa.GetName(),
				Namespace: sa.GetNamespace(),
			},
		},
	}
}

func (b *builder) clusterRoleAndBinding() []client.Object {
	var objs []client.Object
	if b.accessToken.Spec.ClusterPermissions == nil {
		return nil
	}

	clusterRole := b.clusterRole(b.accessToken.Spec.ClusterPermissions.Rules)
	objs = append(objs, clusterRole)

	roleRef := rbacv1.RoleRef{
		APIGroup: rbacv1.GroupName,
		Kind:     "ClusterRole",
		Name:     clusterRole.Name,
	}
	objs = append(objs, b.clusterRoleBinding(roleRef))

	return objs
}

func (b *builder) clusterRole(rules []rbacv1.PolicyRule) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			// NOTE: ClusterRoles are cluster-scoped objects so we qualify the name with the namespace to avoid colliding names
			Name: fmt.Sprintf("%s-%s", b.accessToken.GetName(), b.accessToken.GetNamespace()),
		},
		Rules: rules,
	}
}

func (b *builder) clusterRoleBinding(roleRef rbacv1.RoleRef) *rbacv1.ClusterRoleBinding {
	sa := b.serviceAccount()
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			// NOTE: ClusterRoles are cluster-scoped objects so we qualify the name with the namespace to avoid colliding names
			Name: fmt.Sprintf("%s-%s", b.accessToken.GetName(), b.accessToken.GetNamespace()),
		},
		RoleRef: roleRef,
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      sa.GetName(),
				Namespace: sa.GetNamespace(),
			},
		},
	}
}

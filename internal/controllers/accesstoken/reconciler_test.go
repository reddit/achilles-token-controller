package accesstoken_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/reddit/achilles-token-controller/api/group.example.com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("AccessTokenReconciler", Ordered, func() {
	var accessToken *v1alpha1.AccessToken

	BeforeEach(func() {
		accessToken = &v1alpha1.AccessToken{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foobar",
				Namespace: "default",
			},
			Spec: v1alpha1.AccessTokenSpec{
				NamespacedPermissions: []v1alpha1.NamespacedPermissions{
					{
						Namespace: "default",
						Rules: []rbacv1.PolicyRule{
							{
								APIGroups: []string{""},
								Resources: []string{"configmaps"},
								Verbs:     []string{"*"},
							},
						},
					},
					{
						Namespace: "kube-system",
						Rules: []rbacv1.PolicyRule{
							{
								APIGroups: []string{""},
								Resources: []string{"configmaps"},
								Verbs:     []string{"get", "list", "watch"},
							},
						},
					},
				},
				ClusterPermissions: &v1alpha1.ClusterPermissions{
					Rules: []rbacv1.PolicyRule{
						{
							APIGroups: []string{""},
							Resources: []string{"namespaces"},
							Verbs:     []string{"get", "list", "watch"},
						},
					},
				},
			},
		}

		Expect(c.Create(ctx, accessToken)).To(Succeed())
	})

	It("should reconcile AccessToken objects", func() {
		By("provisioning access token with relevant permissions")

		// Secret
		Eventually(func(g Gomega) {
			expected := &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "default",
				},
			}
			actual := &corev1.Secret{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expected), actual)).To(Succeed())

			g.Expect(actual.Type).To(Equal(corev1.SecretTypeServiceAccountToken))
			g.Expect(actual.Annotations).To(HaveKeyWithValue("kubernetes.io/service-account.name", accessToken.Name))
		}).Should(Succeed())

		// ServiceAccount
		Eventually(func(g Gomega) {
			expected := &corev1.ServiceAccount{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: accessToken.Namespace,
				},
			}
			actual := &corev1.ServiceAccount{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expected), actual)).To(Succeed())
		}).Should(Succeed())

		// Role + RoleBinding for "default" namespace
		Eventually(func(g Gomega) {
			expectedRole := &rbacv1.Role{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "default",
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"configmaps"},
						Verbs:     []string{"*"},
					},
				},
			}
			expectedRoleBinding := &rbacv1.RoleBinding{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "default",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "Role",
					Name:     accessToken.Name,
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      rbacv1.ServiceAccountKind,
						Name:      accessToken.Name,
						Namespace: accessToken.Namespace,
					},
				},
			}

			actualRole := &rbacv1.Role{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedRole), actualRole)).To(Succeed())
			g.Expect(actualRole.Rules).To(Equal(expectedRole.Rules))

			actualRoleBinding := &rbacv1.RoleBinding{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedRoleBinding), actualRoleBinding)).To(Succeed())
			g.Expect(actualRoleBinding.RoleRef).To(Equal(expectedRoleBinding.RoleRef))
			g.Expect(actualRoleBinding.Subjects).To(Equal(expectedRoleBinding.Subjects))
		}).Should(Succeed())

		// Role + RoleBinding for "kube-system" namespace
		Eventually(func(g Gomega) {
			expectedRole := &rbacv1.Role{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "kube-system",
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"configmaps"},
						Verbs:     []string{"get", "list", "watch"},
					},
				},
			}
			expectedRoleBinding := &rbacv1.RoleBinding{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "kube-system",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "Role",
					Name:     accessToken.Name,
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      rbacv1.ServiceAccountKind,
						Name:      accessToken.Name,
						Namespace: accessToken.Namespace,
					},
				},
			}

			actualRole := &rbacv1.Role{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedRole), actualRole)).To(Succeed())
			g.Expect(actualRole.Rules).To(Equal(expectedRole.Rules))

			actualRoleBinding := &rbacv1.RoleBinding{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedRoleBinding), actualRoleBinding)).To(Succeed())
			g.Expect(actualRoleBinding.RoleRef).To(Equal(expectedRoleBinding.RoleRef))
			g.Expect(actualRoleBinding.Subjects).To(Equal(expectedRoleBinding.Subjects))
		}).Should(Succeed())

		// ClusterRole + ClusterRoleBinding
		Eventually(func(g Gomega) {
			expectedClusterRole := &rbacv1.ClusterRole{
				ObjectMeta: v1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", accessToken.Name, accessToken.Namespace),
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"namespaces"},
						Verbs:     []string{"get", "list", "watch"},
					},
				},
			}
			expectedClusterRoleBinding := &rbacv1.ClusterRoleBinding{
				ObjectMeta: v1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", accessToken.Name, accessToken.Namespace),
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "ClusterRole",
					Name:     fmt.Sprintf("%s-%s", accessToken.Name, accessToken.Namespace),
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      rbacv1.ServiceAccountKind,
						Name:      accessToken.Name,
						Namespace: accessToken.Namespace,
					},
				},
			}

			actualClusterRole := &rbacv1.ClusterRole{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedClusterRole), actualClusterRole)).To(Succeed())
			g.Expect(actualClusterRole.Rules).To(Equal(expectedClusterRole.Rules))

			actualClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(expectedClusterRoleBinding), actualClusterRoleBinding)).To(Succeed())
			g.Expect(actualClusterRoleBinding.RoleRef).To(Equal(expectedClusterRoleBinding.RoleRef))
			g.Expect(actualClusterRoleBinding.Subjects).To(Equal(expectedClusterRoleBinding.Subjects))
		}).Should(Succeed())

		By("updating status with ref to secret")

		Eventually(func(g Gomega) {
			actual := &v1alpha1.AccessToken{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: accessToken.Namespace,
				},
			}
			expectedTokenSecretRef := ptr.To(accessToken.Name)

			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(actual), actual)).To(Succeed())
			g.Expect(actual.Status.TokenSecretRef).To(Equal(expectedTokenSecretRef))
		}).Should(Succeed())

		By("cleaning up stale permissions")

		// mutate AccessToken to remove permissions for "kube-system" namespace
		accessToken := accessToken.DeepCopy()
		_, err := controllerutil.CreateOrPatch(ctx, c, accessToken, func() error {
			accessToken.Spec.NamespacedPermissions = []v1alpha1.NamespacedPermissions{
				{
					Namespace: "default",
					Rules: []rbacv1.PolicyRule{
						{
							APIGroups: []string{""},
							Resources: []string{"configmaps"},
							Verbs:     []string{"*"},
						},
					},
				},
			}
			return nil
		})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func(g Gomega) {
			expectedDeletedRole := &rbacv1.Role{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "kube-system",
				},
			}
			expectedDeletedRoleBinding := &rbacv1.RoleBinding{
				ObjectMeta: v1.ObjectMeta{
					Name:      accessToken.Name,
					Namespace: "kube-system",
				},
			}

			g.Expect(errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(expectedDeletedRole), &rbacv1.Role{}))).To(BeTrue())
			g.Expect(errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(expectedDeletedRoleBinding), &rbacv1.RoleBinding{}))).To(BeTrue())
		}).Should(Succeed())

		By("performing finalizer logic and cleaning up cluster scoped resources")
		// delete AccessToken
		Expect(c.Delete(ctx, accessToken)).To(Succeed())

		// assert that ClusterRole and ClusterRoleBinding are deleted
		Eventually(func(g Gomega) {
			expectedDeletedClusterRole := &rbacv1.ClusterRole{
				ObjectMeta: v1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", accessToken.Name, accessToken.Namespace),
				},
			}
			expectedDeletedClusterRoleBinding := &rbacv1.ClusterRoleBinding{
				ObjectMeta: v1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", accessToken.Name, accessToken.Namespace),
				},
			}

			g.Expect(errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(expectedDeletedClusterRole), &rbacv1.ClusterRole{}))).To(BeTrue())
			g.Expect(errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(expectedDeletedClusterRoleBinding), &rbacv1.ClusterRoleBinding{}))).To(BeTrue())
		}).Should(Succeed())
	})
})

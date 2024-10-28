package accesstoken

import (
	"context"

	"github.snooguts.net/reddit/achilles-sdk/pkg/fsm"
	"github.snooguts.net/reddit/achilles-sdk/pkg/fsm/types"
	"github.snooguts.net/reddit/achilles-sdk/pkg/io"
	"github.snooguts.net/reddit/achilles-sdk/pkg/logging"
	"github.snooguts.net/reddit/achilles-sdk/pkg/meta"
	"github.snooguts.net/reddit/achilles-sdk/pkg/sets"
	"github.snooguts.net/reddit/achilles-token-controller/api/group.example.com/v1alpha1"
	"github.snooguts.net/reddit/achilles-token-controller/internal/controlplane"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These kubebuilder markers[0] define the access (RBAC) requirements for the
// controller. They are used to produced appropriate Roles (manifests) that can
// be applied to the cluster. You should add a marker for resource/verb
// combination.
//
// [0]: https://book.kubebuilder.io/reference/markers/rbac.html

// +kubebuilder:rbac:groups=group.example.com,resources=accesstokens;accesstokens/status,verbs=*
// +kubebuilder:rbac:groups="",resources=secrets,verbs=*
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=*
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings,verbs=*
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=*

const (
	controllerName = "AccessToken"
)

type state = types.State[*v1alpha1.AccessToken]

type reconciler struct {
	c      *io.ClientApplicator
	scheme *runtime.Scheme
	log    *zap.SugaredLogger
}

func (r *reconciler) provisionToken() *state {
	return &state{
		Name:      "provision-token",
		Condition: conditionTokenProvisioned,
		Transition: func(
			ctx context.Context,
			accessToken *v1alpha1.AccessToken,
			out *types.OutputSet,
		) (*state, types.Result) {
			builder := newBuilder(accessToken)

			outputs := builder.build()
			for _, o := range outputs {
				var applyOpts []io.ApplyOption

				// NOTE: the achilles-sdk by default adds an owner reference to all objects created by the controller,
				// but we want to avoid this for ClusterRole and ClusterRoleBinding objects since they are cluster-scoped
				// and for any object that is not in the same namespace as the AccessToken

				switch o.(type) {
				case *rbacv1.ClusterRole:
					applyOpts = append(applyOpts, io.WithoutOwnerRefs())
				case *rbacv1.ClusterRoleBinding:
					applyOpts = append(applyOpts, io.WithoutOwnerRefs())
				default:
					if o.GetNamespace() != accessToken.GetNamespace() {
						applyOpts = append(applyOpts, io.WithoutOwnerRefs())
					}
				}

				out.Apply(o, applyOpts...)
			}

			accessToken.Status.TokenSecretRef = ptr.To(builder.secret().Name)

			return r.deleteStalePermissions(outputs), types.DoneResult()
		},
	}
}

func (r *reconciler) deleteStalePermissions(desiredObjs []client.Object) *state {
	return &state{
		Name:      "delete-stale-permissions",
		Condition: conditionStalePermissionsRemoved,
		Transition: func(
			ctx context.Context,
			accessToken *v1alpha1.AccessToken,
			out *types.OutputSet,
		) (*state, types.Result) {
			desired := sets.NewObjectSet(r.scheme, desiredObjs...)
			actual := sets.NewObjectSet(r.scheme)

			// get all existing managed resources
			for _, ref := range accessToken.Status.ResourceRefs {
				// get the object from the cluster
				obj, err := meta.NewObjectForGVK(r.scheme, ref.GroupVersionKind())
				if err != nil {
					return nil, types.ErrorResultf("constructing new %T %s: %s", obj, client.ObjectKeyFromObject(obj), err)
				}
				obj.SetName(ref.Name)
				obj.SetNamespace(ref.Namespace)
				if err := r.c.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
					if errors.IsNotFound(err) {
						// warn for managed resource that wasn't explicitly deleted by the controller, but is deleted on the kube-apiserver
						r.log.Warnf("managed resource %T %s is not found", obj, client.ObjectKeyFromObject(obj))
						continue
					}
					return nil, types.ErrorResultf("getting managed object %T %s: %s", obj, client.ObjectKeyFromObject(obj), err)
				}

				actual.Insert(obj)
			}

			// delete stale permissions
			for _, staleObj := range actual.Difference(desired).List() {
				out.Delete(staleObj)
			}

			return nil, types.DoneResult()
		},
	}
}

func SetupController(
	ctx context.Context,
	cpCtx controlplane.Context,
	mgr ctrl.Manager,
	rl workqueue.RateLimiter,
	c *io.ClientApplicator,
) error {
	_, log, err := logging.ControllerCtx(ctx, controllerName)
	if err != nil {
		return err
	}

	r := &reconciler{
		c:      c,
		scheme: mgr.GetScheme(),
		log:    log,
	}

	builder := fsm.NewBuilder(
		&v1alpha1.AccessToken{},
		r.provisionToken(),
		mgr.GetScheme(),
	).Manages(
		corev1.SchemeGroupVersion.WithKind("Secret"),
		corev1.SchemeGroupVersion.WithKind("ServiceAccount"),
		rbacv1.SchemeGroupVersion.WithKind("Role"),
		rbacv1.SchemeGroupVersion.WithKind("RoleBinding"),
		rbacv1.SchemeGroupVersion.WithKind("ClusterRole"),
		rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding"),
	).WithFinalizerState(
		// NOTE: we can't rely on native Kubernetes GC to delete cluster scoped resources (ClusterRole, ClusterRoleBinding)
		// or cross-namespace resources (Roles, RoleBindings) so we need to handle this ourselves
		r.deleteStalePermissions(nil),
	)

	return builder.Build()(mgr, log, rl, cpCtx.Metrics)
}

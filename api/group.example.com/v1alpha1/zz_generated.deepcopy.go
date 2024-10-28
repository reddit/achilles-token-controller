//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.snooguts.net/reddit/achilles-sdk-api/api"
	v1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessToken) DeepCopyInto(out *AccessToken) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessToken.
func (in *AccessToken) DeepCopy() *AccessToken {
	if in == nil {
		return nil
	}
	out := new(AccessToken)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AccessToken) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessTokenList) DeepCopyInto(out *AccessTokenList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AccessToken, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessTokenList.
func (in *AccessTokenList) DeepCopy() *AccessTokenList {
	if in == nil {
		return nil
	}
	out := new(AccessTokenList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AccessTokenList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessTokenSpec) DeepCopyInto(out *AccessTokenSpec) {
	*out = *in
	if in.NamespacedPermissions != nil {
		in, out := &in.NamespacedPermissions, &out.NamespacedPermissions
		*out = make([]NamespacedPermissions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ClusterPermissions != nil {
		in, out := &in.ClusterPermissions, &out.ClusterPermissions
		*out = new(ClusterPermissions)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessTokenSpec.
func (in *AccessTokenSpec) DeepCopy() *AccessTokenSpec {
	if in == nil {
		return nil
	}
	out := new(AccessTokenSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessTokenStatus) DeepCopyInto(out *AccessTokenStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.ResourceRefs != nil {
		in, out := &in.ResourceRefs, &out.ResourceRefs
		*out = make([]api.TypedObjectRef, len(*in))
		copy(*out, *in)
	}
	if in.TokenSecretRef != nil {
		in, out := &in.TokenSecretRef, &out.TokenSecretRef
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessTokenStatus.
func (in *AccessTokenStatus) DeepCopy() *AccessTokenStatus {
	if in == nil {
		return nil
	}
	out := new(AccessTokenStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPermissions) DeepCopyInto(out *ClusterPermissions) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]v1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPermissions.
func (in *ClusterPermissions) DeepCopy() *ClusterPermissions {
	if in == nil {
		return nil
	}
	out := new(ClusterPermissions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedPermissions) DeepCopyInto(out *NamespacedPermissions) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]v1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedPermissions.
func (in *NamespacedPermissions) DeepCopy() *NamespacedPermissions {
	if in == nil {
		return nil
	}
	out := new(NamespacedPermissions)
	in.DeepCopyInto(out)
	return out
}

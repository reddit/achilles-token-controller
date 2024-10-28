# achilles-token-controller

This is an example Achilles SDK based controller showcasing SDK basics.
It implements the `AccessToken` CRD, which allows creating a Kubernetes bearer token with
specified permissions.

## Running the controller

1. Clone the `achilles-token-controller`.

    ```
    git clone git@github.snooguts.net:reddit/achilles-token-controller.git
    ```

1. Ensure you have [k3d](https://k3d.io/v5.7.4/#installation) installed.

1. Deploy a local cluster with k3d.

    ```sh
    k3d cluster create orch
    ```

1. Verify the above command updated your `kubecontext` to the k3d cluster.

    ```sh
    kubectl config current-context
    ```

   The output should be:

    ```txt
    k3d-orch
    ```
1. Build the controller image.

    ```sh
    make docker
    ```

1. Load the controller image into the k3d cluster

   ```sh
   k3d image import achilles-token-controller:latest -c orch
   ```

1. Open `manifests/base/manager.yaml` and replace `image: REPLACE-ME` with `image: achilles-token-controller:latest`.
   If this file doesn't exist, run `make generate`.
1. Create the namespace for the controller
   ```sh
   kubectl create namespace achilles-system
   ```
1. Deploy the controller.
    ```sh
    kubectl apply -f manifests/base/manager.yaml
    ```
1. Test the controller with this example AccessToken.
   ```yaml
   apiVersion: group.example.com/v1alpha1
   kind: AccessToken
   metadata:
     name: test
     namespace: default
   spec:
     namespacedPermissions:
     - namespace: default
       rules:
       - apiGroups: [""]
         resources: ["configmaps"]
         verbs:     ["*"]
     - namespace: kube-system
       rules:
       - apiGroups: [""]
         resources: ["configmaps"]
         verbs:     ["get", "list", "watch"]
     clusterPermissions:
       rules:
       - apiGroups: [""]
         resources: ["namespaces"]
         verbs:     ["get", "list", "watch"]
    ```
1. Check that the AccessToken was processed successfully
   ```sh
   kubectl get accesstoken test -n default -oyaml
   ```

   You should see the following status condition, indicating that the object was instantiated successfully.

   ```yaml
    status:
      conditions:
      - lastTransitionTime: "2024-10-24T17:33:35Z"
        message: All conditions successful.
        observedGeneration: 1
        reason: ConditionsSuccessful
        status: "True"
        type: Ready
    ```
   You'll also see that it provisioned a deploy token as a secret, whose name is under `status.tokenSecretRef`.

1. As a bonus, we can use `kubectl auth can-i` ([docs here](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_auth/kubectl_auth_can-i/))
   check that the deploy token in fact has the permissions that we declared for it.
   We first need to locate the Service Account that the AccessToken was created for, which can be found under `status.resourceRefs`
   with `kind: ServiceAccount`.

    ```sh
    kubectl auth can-i --as=system:serviceaccount:default:test create configmaps -n default # should report yes
    kubectl auth can-i --as=system:serviceaccount:default:test create configmaps -n kube-system # should report no
    kubectl auth can-i --as=system:serviceaccount:default:test list configmaps -n kube-system # should report yes
    kubectl auth can-i --as=system:serviceaccount:default:test create namespaces # should report no
    kubectl auth can-i --as=system:serviceaccount:default:test list namespaces # should report yes
    ```

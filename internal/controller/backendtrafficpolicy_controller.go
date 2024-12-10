/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/utils/datastore"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

// BackendTrafficPolicyReconciler reconciles a BackendTrafficPolicy object
type BackendTrafficPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gateway.envoyproxy.io,resources=backendtrafficpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.envoyproxy.io,resources=backendtrafficpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.envoyproxy.io,resources=backendtrafficpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BackendTrafficPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *BackendTrafficPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic
	log.Info("Reconciling BackendTrafficPolicy")

	// read BTP instance from btps.json file and, fetch each corresponding instance from controller-runtime

	// get dataStoreFilePath
	fileStore := datastore.NewFileStore(log.WithName("FileStore"))
	dataStoreFilePath, err := fileStore.GetDataStoreFilePath()
	if err != nil {
		log.Error(err, "Error getting dataStoreFilePath")
		return ctrl.Result{}, err
	}

	// read BTPs from file
	btps, err := fileStore.ReadBTPsFromFile(dataStoreFilePath)
	if err != nil {
		log.Error(err, "Error reading BTPs from file")
		return ctrl.Result{}, err
	}

	log.Info("BTPs from file", "BTPs", btps)
	// fetch each BTP instance from controller-runtime
	for _, btp := range btps {
		name := btp.Metadata.Name
		namespace := btp.Metadata.Namespace
		// fetch BackendTrafficPolicy instance
		btpInstance := &egv1a1.BackendTrafficPolicy{}
		// expected BTP instance
		expectedBackendTrafficPolicy := r.expectedBackendTrafficPolicy(&btp)
		// expectedBackendTrafficPolicy := r.defineBackendTrafficPolicy(&btp)
		if err := r.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, btpInstance); err != nil {
			if kerrors.IsNotFound(err) {
				log.Info("BTP instance not found", "BTP", btp)

				// set the owner reference
				// ownerRef := metav1.NewControllerRef(expectedBackendTrafficPolicy, egv1a1.SchemeGroupVersion.WithKind(egv1a1.KindBackendTrafficPolicy))
				// btpInstance.SetOwnerReferences([]metav1.OwnerReference{*ownerRef})

				// set controller reference
				// if err := ctrl.SetControllerReference(expectedBackendTrafficPolicy, btpInstance, r.Scheme); err != nil {
				// 	log.Error(err, "Failed to set controller reference")

				// 	// update the status of the BTP instance
				// 	// btp.Status = gwapiv1a2.PolicyStatus{
				// 	// 	State:  gwapiv1a2.Rejected,
				// 	// 	Reason: "Failed to set controller reference",
				// 	// }

				// 	if err := r.Status().Update(ctx, expectedBackendTrafficPolicy); err != nil {
				// 		log.Error(err, "Failed to update BTP status")
				// 		return ctrl.Result{}, err
				// 	}

				// 	return ctrl.Result{}, err

				// }

				// create BTP instance
				log.Info("Creating new BTP instance",
					"BackendTrafficPolicy.Namespace", expectedBackendTrafficPolicy.Namespace,
					"BackendTrafficPolicy.Name", expectedBackendTrafficPolicy.Name)

				if err := r.Create(ctx, expectedBackendTrafficPolicy); err != nil {
					log.Error(err, "Failed to create new BTP instance",
						"BackendTrafficPolicy.Namespace", expectedBackendTrafficPolicy.Namespace,
						"BackendTrafficPolicy.Name", expectedBackendTrafficPolicy.Name)

					return ctrl.Result{}, err
				}

				// btp created successfully
				log.Info("BTP instance created successfully",
					"BackendTrafficPolicy.Namespace", expectedBackendTrafficPolicy.Namespace,
					"BackendTrafficPolicy.Name", expectedBackendTrafficPolicy.Name)

				// update the status of the BTP instance

				return ctrl.Result{}, nil

			}

			// error reading the BTP instance
			log.Error(err, "Failed to get BTP instance")
			return ctrl.Result{}, err
		}

		// BTP instance found

		// verify if the BTP instance is same as expected
		currentBTPSpec := btpInstance.Spec
		expectedBTPSpec := expectedBackendTrafficPolicy.Spec

		if r.compareBackendTrafficPolicy(&expectedBTPSpec, &currentBTPSpec) {
			log.Info("BTP instance is not same as expected",
				"BackendTrafficPolicy.Namespace", btpInstance.Namespace,
				"BackendTrafficPolicy.Name", btpInstance.Name)

			// update the BTP instance
			btpInstance.Spec = expectedBTPSpec
			if err := r.Update(ctx, btpInstance); err != nil {
				log.Error(err, "Failed to update BTP instance",
					"BackendTrafficPolicy.Namespace", btpInstance.Namespace,
					"BackendTrafficPolicy.Name", btpInstance.Name)

				return ctrl.Result{}, err
			}

			// BTP instance updated successfully
			log.Info("BTP instance updated successfully",
				"BackendTrafficPolicy.Namespace", btpInstance.Namespace,
				"BackendTrafficPolicy.Name", btpInstance.Name)

			return ctrl.Result{}, nil
		}

		// BTP instance is same as expected
		log.Info("BTP instance is same as expected, no change required",
			"BackendTrafficPolicy.Namespace", btpInstance.Namespace,
			"BackendTrafficPolicy.Name", btpInstance.Name)

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackendTrafficPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&egv1a1.BackendTrafficPolicy{}).
		Complete(r)
}

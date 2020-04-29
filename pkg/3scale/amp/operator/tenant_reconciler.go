package operator

import (
	"context"
	"fmt"

	capabilitiesv1alpha1 "github.com/3scale/3scale-operator/pkg/apis/capabilities/v1alpha1"
	"github.com/3scale/3scale-operator/pkg/helper"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("tenant_reconciler")

// type TenantReconciler interface {
// 	IsUpdateNeeded(desired, existing *capabilitiesv1alpha1.Tenant) bool
// }

type TenantReconciler struct {
	BaseAPIManagerLogicReconciler
}

type TenantBaseReconciler struct {
	BaseAPIManagerLogicReconciler
	reconciler TenantReconciler
}

func NewTenantReconciler(baseAPIManagerLogicReconciler BaseAPIManagerLogicReconciler) *TenantReconciler {
	return &TenantReconciler{
		BaseAPIManagerLogicReconciler: baseAPIManagerLogicReconciler,
	}
}

func NewTenantBaseReconciler(baseAPIManagerLogicReconciler BaseAPIManagerLogicReconciler, reconciler TenantReconciler) *TenantBaseReconciler {
	return &TenantBaseReconciler{
		BaseAPIManagerLogicReconciler: baseAPIManagerLogicReconciler,
		reconciler:                    reconciler,
	}
}

func (r *TenantBaseReconciler) Reconcile(desired *capabilitiesv1alpha1.Tenant) error {
	reqLogger := log.WithValues("Desired.Namespace", desired.GetNamespace(), "Desired.Name", desired.GetName())
	reqLogger.Info("Reconciling Tenant")

	objectInfo := ObjectInfo(desired)
	namespace := desired.GetNamespace()
	if namespace == "" {
		namespace = r.apiManager.GetNamespace()
	}
	existing := &capabilitiesv1alpha1.Tenant{}
	err := r.Client().Get(
		context.TODO(),
		types.NamespacedName{Name: desired.Name, Namespace: namespace},
		existing)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Tenant resource not found")
			createErr := r.createResource(desired)
			if createErr != nil {
				r.Logger().Error(createErr, fmt.Sprintf("Error creating object %s. Requeuing request...", objectInfo))
				return createErr
			}
			return nil
		}
		return err
	}

	update, err := r.isUpdateNeeded(desired, existing)
	if err != nil {
		return err
	}

	if update {
		return r.updateResource(existing)
	}

	return nil
}

func (r *TenantBaseReconciler) isUpdateNeeded(desired, existing *capabilitiesv1alpha1.Tenant) (bool, error) {
	updated := helper.EnsureObjectMeta(&existing.ObjectMeta, &desired.ObjectMeta)

	updatedTmp, err := r.ensureOwnerReference(existing)
	if err != nil {
		return false, nil
	}

	updated = updated || updatedTmp

	// // updatedTmp = r.reconciler.IsUpdateNeeded(desired, existing)
	// updated = updated || updatedTmp

	return updated, nil
}

// reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
// reqLogger.Info("Reconciling Tenant")

// // Fetch the Tenant instance
// tenantR := &apiv1alpha1.Tenant{}
// err := r.client.Get(context.TODO(), request.NamespacedName, tenantR)
// if err != nil {
// 	if errors.IsNotFound(err) {
// 		// Request object not found, could have been deleted after reconcile request.
// 		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
// 		// Return and don't requeue
// 		reqLogger.Info("Tenant resource not found")
// 		return reconcile.Result{}, nil
// 	}
// 	// Error reading the object - requeue the request.
// 	return reconcile.Result{}, err
// }

// changed := tenantR.SetDefaults()
// if changed {
// 	err = r.client.Update(context.TODO(), tenantR)
// 	if err != nil {
// 		return reconcile.Result{}, err
// 	}
// 	reqLogger.Info("Tenant resource updated with defaults")
// 	// Expect for re-trigger
// 	return reconcile.Result{}, nil
// }

// masterAccessToken, err := FetchMasterCredentials(r.client, tenantR)
// if err != nil {
// 	log.Error(err, "Error fetching master credentials secret")
// 	// Error reading the object - requeue the request.
// 	return reconcile.Result{}, err
// }

// portaClient, err := helper.PortaClientFromURLString(tenantR.Spec.SystemMasterUrl, masterAccessToken)
// if err != nil {
// 	log.Error(err, "Error creating porta client object")
// 	// Error reading the object - requeue the request.
// 	return reconcile.Result{}, err
// }

// internalReconciler := NewInternalReconciler(r.client, tenantR, portaClient, reqLogger)
// err = internalReconciler.Run()
// if err != nil {
// 	log.Error(err, "Error in tenant reconciliation")
// 	// Error reading the object - requeue the request.
// 	return reconcile.Result{}, err
// }

// reqLogger.Info("Tenant reconciled successfully")
// return reconcile.Result{}, nil
// }

// // FetchMasterCredentials get secret using k8s client
// func FetchMasterCredentials(k8sClient client.Client, tenantR *apiv1alpha1.Tenant) (string, error) {
// masterCredentialsSecret := &v1.Secret{}

// err := k8sClient.Get(context.TODO(),
// 	types.NamespacedName{
// 		Name:      tenantR.Spec.MasterCredentialsRef.Name,
// 		Namespace: tenantR.Spec.MasterCredentialsRef.Namespace,
// 	},
// 	masterCredentialsSecret)

// if err != nil {
// 	return "", err
// }

// masterAccessTokenByteArray, ok := masterCredentialsSecret.Data[component.SystemSecretSystemSeedMasterAccessTokenFieldName]
// if !ok {
// 	return "", fmt.Errorf("Key not found in master secret (ns: %s, name: %s) key: %s",
// 		tenantR.Spec.MasterCredentialsRef.Namespace, tenantR.Spec.MasterCredentialsRef.Name,
// 		component.SystemSecretSystemSeedMasterAccessTokenFieldName)
// }

// return bytes.NewBuffer(masterAccessTokenByteArray).String(), nil

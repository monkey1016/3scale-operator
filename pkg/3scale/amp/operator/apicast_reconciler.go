package operator

import (
	"context"

	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	appsv1alpha1 "github.com/3scale/3scale-operator/pkg/apis/apps/v1alpha1"
	capabilitiesv1alpha1 "github.com/3scale/3scale-operator/pkg/apis/capabilities/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ApicastEnvCMReconciler struct {
}

func NewApicastEnvCMReconciler() *ApicastEnvCMReconciler {
	return &ApicastEnvCMReconciler{}
}

func (r *ApicastEnvCMReconciler) IsUpdateNeeded(desiredCM, existingCM *v1.ConfigMap) bool {
	update := false

	//	Check APICAST_MANAGEMENT_API
	fieldUpdated := ConfigMapReconcileField(desiredCM, existingCM, "APICAST_MANAGEMENT_API")
	update = update || fieldUpdated

	//	Check OPENSSL_VERIFY
	fieldUpdated = ConfigMapReconcileField(desiredCM, existingCM, "OPENSSL_VERIFY")
	update = update || fieldUpdated

	//	Check APICAST_RESPONSE_CODES
	fieldUpdated = ConfigMapReconcileField(desiredCM, existingCM, "APICAST_RESPONSE_CODES")
	update = update || fieldUpdated

	return update
}

type ApicastStagingDCReconciler struct {
	BaseAPIManagerLogicReconciler
}

func NewApicastDCReconciler(baseAPIManagerLogicReconciler BaseAPIManagerLogicReconciler) *ApicastStagingDCReconciler {
	return &ApicastStagingDCReconciler{
		BaseAPIManagerLogicReconciler: baseAPIManagerLogicReconciler,
	}
}

func (r *ApicastStagingDCReconciler) IsUpdateNeeded(desired, existing *appsv1.DeploymentConfig) bool {
	update := false

	tmpUpdate := DeploymentConfigReconcileReplicas(desired, existing, r.Logger())
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigReconcileContainerResources(desired, existing, r.Logger())
	update = update || tmpUpdate

	return update
}

type ApicastReconciler struct {
	BaseAPIManagerLogicReconciler
	apicastSpec *appsv1alpha1.ApicastSpec
}

// blank assignment to verify that BaseReconciler implements reconcile.Reconciler
var _ LogicReconciler = &ApicastReconciler{}

func NewApicastReconciler(baseAPIManagerLogicReconciler BaseAPIManagerLogicReconciler, apicast *appsv1alpha1.ApicastSpec) ApicastReconciler {
	return ApicastReconciler{
		BaseAPIManagerLogicReconciler: baseAPIManagerLogicReconciler,
		apicastSpec:                   apicast,
	}
}

func (r *ApicastReconciler) Reconcile() (reconcile.Result, error) {
	apicast, err := r.apicast()
	if err != nil {
		return reconcile.Result{}, err
	}

	var portalEnpointSecret *string
	if *(r.apicastSpec.CreateTenant) {
		// Create a tenant for this
		environmentTenant := apicast.EnvironmentTenant(&r.apiManager.Namespace)
		err = r.reconcileTenant(environmentTenant)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Get the THREESCALE_PORTAL_ENDPOINT for the newly created tenant
		existingSecret := &v1.Secret{}
		existingSecret.Name = environmentTenant.Spec.TenantSecretRef.Name
		existingSecret.Namespace = environmentTenant.Spec.TenantSecretRef.Namespace

		err := r.Client().Get(
			context.TODO(),
			types.NamespacedName{Name: existingSecret.Name, Namespace: existingSecret.Namespace},
			existingSecret)
		if err != nil {
			return reconcile.Result{}, err
		}
		portalEnpointSecret = &existingSecret.Name
	}
	err = r.reconcileDeploymentConfig(apicast.DeploymentConfig(&r.apiManager.Namespace, portalEnpointSecret))
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileService(apicast.Service())
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileEnvironmentConfigMap(apicast.EnvironmentConfigMap())
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ApicastReconciler) apicast() (*component.Apicast, error) {
	optsProvider := OperatorApicastOptionsProvider{APIManagerSpec: &r.apiManager.Spec, Namespace: r.apiManager.Namespace, Client: r.Client(), ApicastSpec: r.apicastSpec}
	opts, err := optsProvider.GetApicastOptions()
	if err != nil {
		return nil, err
	}
	return component.NewApicast(opts), nil
}

// func (r *ApicastReconciler) reconcileStagingDeploymentConfig(desiredDeploymentConfig *appsv1.DeploymentConfig) error {
// 	reconciler := NewDeploymentConfigBaseReconciler(r.BaseAPIManagerLogicReconciler, NewApicastDCReconciler(r.BaseAPIManagerLogicReconciler))
// 	return reconciler.Reconcile(desiredDeploymentConfig)
// }

func (r *ApicastReconciler) reconcileDeploymentConfig(desiredDeploymentConfig *appsv1.DeploymentConfig) error {
	reconciler := NewDeploymentConfigBaseReconciler(r.BaseAPIManagerLogicReconciler, NewApicastDCReconciler(r.BaseAPIManagerLogicReconciler))
	return reconciler.Reconcile(desiredDeploymentConfig)
}

func (r *ApicastReconciler) reconcileTenant(desiredTentant *capabilitiesv1alpha1.Tenant) error {
	reconciler := NewTenantBaseReconciler(r.BaseAPIManagerLogicReconciler, *(NewTenantReconciler(r.BaseAPIManagerLogicReconciler)))
	return reconciler.Reconcile(desiredTentant)
}

// func (r *ApicastReconciler) reconcileProductionDeploymentConfig(desiredDeploymentConfig *appsv1.DeploymentConfig) error {
// 	reconciler := NewDeploymentConfigBaseReconciler(r.BaseAPIManagerLogicReconciler, NewApicastDCReconciler(r.BaseAPIManagerLogicReconciler))
// 	return reconciler.Reconcile(desiredDeploymentConfig)
// }

// func (r *ApicastReconciler) reconcileStagingService(desiredService *v1.Service) error {
// 	reconciler := NewServiceBaseReconciler(r.BaseAPIManagerLogicReconciler, NewCreateOnlySvcReconciler())
// 	return reconciler.Reconcile(desiredService)
// }

// func (r *ApicastReconciler) reconcileProductionService(desiredService *v1.Service) error {
// 	reconciler := NewServiceBaseReconciler(r.BaseAPIManagerLogicReconciler, NewCreateOnlySvcReconciler())
// 	return reconciler.Reconcile(desiredService)
// }

func (r *ApicastReconciler) reconcileService(desiredService *v1.Service) error {
	reconciler := NewServiceBaseReconciler(r.BaseAPIManagerLogicReconciler, NewCreateOnlySvcReconciler())
	return reconciler.Reconcile(desiredService)
}

func (r *ApicastReconciler) reconcileEnvironmentConfigMap(desiredConfigMap *v1.ConfigMap) error {
	reconciler := NewConfigMapBaseReconciler(r.BaseAPIManagerLogicReconciler, NewApicastEnvCMReconciler())
	return reconciler.Reconcile(desiredConfigMap)
}

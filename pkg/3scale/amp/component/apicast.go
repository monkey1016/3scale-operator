package component

import (
	"github.com/3scale/3scale-operator/pkg/common"

	capabilitiesv1alpha1 "github.com/3scale/3scale-operator/pkg/apis/capabilities/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Apicast struct {
	Options *ApicastOptions
}

func NewApicast(options *ApicastOptions) *Apicast {
	return &Apicast{Options: options}
}

func (apicast *Apicast) Objects() []common.KubernetesObject {
	deploymentConfig := apicast.DeploymentConfig(nil, nil)
	// stagingDeploymentConfig := apicast.StagingDeploymentConfig()
	// productionDeploymentConfig := apicast.ProductionDeploymentConfig()
	service := apicast.Service()
	// stagingService := apicast.StagingService()
	// productionService := apicast.ProductionService()
	environmentConfigMap := apicast.EnvironmentConfigMap()

	objects := []common.KubernetesObject{
		deploymentConfig,
		service,
		environmentConfigMap,
	}
	return objects
}

func (apicast *Apicast) buildApicastCommonEnv(portalEndpointSecret *string) []v1.EnvVar {
	defaultPortalEndpointSecret := "system-master-apicast"
	if portalEndpointSecret == nil {
		portalEndpointSecret = &defaultPortalEndpointSecret
	}
	return []v1.EnvVar{
		envVarFromSecret("THREESCALE_PORTAL_ENDPOINT", *portalEndpointSecret, "PROXY_CONFIGS_ENDPOINT"),
		envVarFromSecret("BACKEND_ENDPOINT_OVERRIDE", "backend-listener", "service_endpoint"),
		envVarFromConfigMap("APICAST_MANAGEMENT_API", "apicast-environment", "APICAST_MANAGEMENT_API"),
		envVarFromConfigMap("OPENSSL_VERIFY", "apicast-environment", "OPENSSL_VERIFY"),
		envVarFromConfigMap("APICAST_RESPONSE_CODES", "apicast-environment", "APICAST_RESPONSE_CODES"),
	}
}

func (apicast *Apicast) buildApicastStagingEnv() []v1.EnvVar {
	result := []v1.EnvVar{}
	result = append(result, apicast.buildApicastCommonEnv(nil)...)
	result = append(result,
		envVarFromValue("APICAST_CONFIGURATION_LOADER", "lazy"),
		envVarFromValue("APICAST_CONFIGURATION_CACHE", "0"),
		envVarFromValue("THREESCALE_DEPLOYMENT_ENV", "staging"),
	)
	return result
}

func (apicast *Apicast) buildApicastProductionEnv(portalEndpointSecret *string) []v1.EnvVar {
	result := []v1.EnvVar{}
	result = append(result, apicast.buildApicastCommonEnv(portalEndpointSecret)...)
	result = append(result,
		envVarFromValue("APICAST_CONFIGURATION_LOADER", "lazy"),
		envVarFromValue("APICAST_CONFIGURATION_CACHE", "0"),
		envVarFromValue("THREESCALE_DEPLOYMENT_ENV", "production"),
	)
	return result
}

func (apicast *Apicast) EnvironmentConfigMap() *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apicast-environment",
			Namespace: *apicast.Options.namespace,
			Labels:    map[string]string{"threescale_component": "apicast", "app": apicast.Options.appLabel},
		},
		Data: map[string]string{
			"APICAST_MANAGEMENT_API": apicast.Options.managementAPI,
			"OPENSSL_VERIFY":         apicast.Options.openSSLVerify,
			"APICAST_RESPONSE_CODES": apicast.Options.responseCodes,
		},
	}
}

func (apicast *Apicast) EnvironmentTenant(systemNamespace *string) *capabilitiesv1alpha1.Tenant {
	return &capabilitiesv1alpha1.Tenant{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capabilitiesv1alpha1.SchemeGroupVersion.String(),
			Kind:       capabilitiesv1alpha1.TenantKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: *apicast.Options.namespace,
			Name:      *apicast.Options.environment,
			Labels: map[string]string{
				"app":                          apicast.Options.appLabel,
				"threescale_component":         "apicast",
				"threescale_component_element": *apicast.Options.environment,
			},
		},
		Spec: capabilitiesv1alpha1.TenantSpec{
			Username:         "admin",
			SystemMasterUrl:  *(getSystemMasterUrl(systemNamespace)),
			Email:            "admin@3scale-operator.com",
			OrganizationName: *apicast.Options.environment,
			MasterCredentialsRef: v1.SecretReference{
				Name:      "system-seed",
				Namespace: *systemNamespace,
			},
			PasswordCredentialsRef: v1.SecretReference{
				Name:      *apicast.Options.environment,
				Namespace: *apicast.Options.namespace,
			},
			TenantSecretRef: v1.SecretReference{
				Name:      "tenant-" + *apicast.Options.environment,
				Namespace: *apicast.Options.namespace,
			},
		},
	}
}

func getSystemMasterUrl(systemNamespace *string) *string {
	systemMasterURL := "http://system-master:3000/status"
	if systemNamespace != nil {
		systemMasterURL = "http://system-master." + *systemNamespace + ".svc.cluster.local:3000/status"
	}
	return &systemMasterURL
}

func (apicast *Apicast) DeploymentConfig(systemNamespace *string, portalEndpointSecret *string) *appsv1.DeploymentConfig {
	systemMasterUrl := *(getSystemMasterUrl(systemNamespace))
	return &appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps.openshift.io/v1", Kind: "DeploymentConfig"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-" + *apicast.Options.environment,
			Labels: map[string]string{
				"app":                          apicast.Options.appLabel,
				"threescale_component":         "apicast",
				"threescale_component_element": *apicast.Options.environment,
			},
			Namespace: *apicast.Options.namespace,
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas: *apicast.Options.replicas,
			Selector: map[string]string{
				"deploymentConfig": "apicast-" + *apicast.Options.environment,
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingParams: &appsv1.RollingDeploymentStrategyParams{
					IntervalSeconds: &[]int64{1}[0],
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(intstr.String),
						StrVal: "25%",
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(intstr.String),
						StrVal: "25%",
					},
					TimeoutSeconds:      &[]int64{1800}[0],
					UpdatePeriodSeconds: &[]int64{1}[0],
				},
				Type: appsv1.DeploymentStrategyTypeRolling,
			},
			Triggers: appsv1.DeploymentTriggerPolicies{
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnConfigChange,
				},
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							"system-master-svc",
							"apicast-" + *apicast.Options.environment,
						},
						From: v1.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "amp-apicast:latest",
						},
					},
				},
			},
			Template: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deploymentConfig":             "apicast-" + *apicast.Options.environment,
						"app":                          apicast.Options.appLabel,
						"threescale_component":         "apicast",
						"threescale_component_element": *apicast.Options.environment,
					},
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/port":   "9421",
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "amp",
					InitContainers: []v1.Container{
						v1.Container{
							Name:    "system-master-svc",
							Image:   "amp-apicast:latest",
							Command: []string{"sh", "-c", "until $(curl --output /dev/null --silent --fail --head " + systemMasterUrl + "); do sleep $SLEEP_SECONDS; done"},
							Env: []v1.EnvVar{
								v1.EnvVar{
									Name:  "SLEEP_SECONDS",
									Value: "1",
								},
							},
						},
					},
					Containers: []v1.Container{
						v1.Container{
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 8080,
									Protocol:      v1.ProtocolTCP,
								},
								v1.ContainerPort{
									ContainerPort: 8090,
									Protocol:      v1.ProtocolTCP,
								},
								v1.ContainerPort{
									ContainerPort: 9421,
									Protocol:      v1.ProtocolTCP,
									Name:          "metrics",
								},
							},
							Env:             apicast.buildApicastProductionEnv(portalEndpointSecret),
							Image:           "amp-apicast:latest",
							ImagePullPolicy: v1.PullIfNotPresent,
							Name:            "apicast-" + *apicast.Options.environment,
							Resources:       *apicast.Options.resourceRequirements,
							LivenessProbe: &v1.Probe{
								Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
									Path: "/status/live",
									Port: intstr.FromInt(8090),
								}},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      5,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
									Path: "/status/ready",
									Port: intstr.FromInt(8090),
								}},
								InitialDelaySeconds: 15,
								TimeoutSeconds:      5,
								PeriodSeconds:       30,
							},
						},
					},
				},
			},
		},
	}
}

func (apicast *Apicast) Service() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-" + *apicast.Options.environment,
			Labels: map[string]string{
				"app":                          apicast.Options.appLabel,
				"threescale_component":         "apicast",
				"threescale_component_element": *apicast.Options.environment,
			},
			Namespace: *apicast.Options.namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       "gateway",
					Protocol:   v1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
				v1.ServicePort{
					Name:       "management",
					Protocol:   v1.ProtocolTCP,
					Port:       8090,
					TargetPort: intstr.FromInt(8090),
				},
			},
			Selector: map[string]string{"deploymentConfig": "apicast-" + *apicast.Options.environment},
		},
	}
}

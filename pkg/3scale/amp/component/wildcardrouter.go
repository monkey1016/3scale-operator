package component

import (
	appsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type WildcardRouter struct {
	options []string
}

func NewWildcardRouter(options []string) *WildcardRouter {
	wr := &WildcardRouter{
		options: options,
	}
	return wr
}

func (wr *WildcardRouter) AssembleIntoTemplate(template *templatev1.Template, otherComponents []Component) {
	wr.buildParameters(template)
	wr.buildObjects(template)
}

func (wr *WildcardRouter) PostProcess(template *templatev1.Template, otherComponents []Component) {

}

func (wr *WildcardRouter) buildParameters(template *templatev1.Template) {
	parameters := []templatev1.Parameter{
		templatev1.Parameter{
			Name:        "WILDCARD_DOMAIN",
			Description: "Root domain for the wildcard routes. Eg. example.com will generate 3scale-admin.example.com.",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "WILDCARD_POLICY",
			Description: "Use \"Subdomain\" to create a wildcard route for apicast wildcard router",
			Value:       "None",
			Required:    true,
		},
	}
	template.Parameters = append(template.Parameters, parameters...)
}

func (wr *WildcardRouter) buildObjects(template *templatev1.Template) {
	wildcardRouterDeploymentConfig := wr.buildWildcardRouterDeploymentConfig()
	wildcardRouterService := wr.buildWildcardRouterService()
	wildcardRouterRoute := wr.buildWildcardRouterRoute()

	objects := []runtime.RawExtension{
		runtime.RawExtension{Object: wildcardRouterDeploymentConfig},
		runtime.RawExtension{Object: wildcardRouterService},
		runtime.RawExtension{Object: wildcardRouterRoute},
	}
	template.Objects = append(template.Objects, objects...)
}

func (wr *WildcardRouter) buildWildcardRouterRoute() *routev1.Route {
	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "apicast-wildcard-router",
			Labels: map[string]string{"app": "${APP_LABEL}", "3scale.component": "apicast", "3scale.component-element": "wildcard-router"},
		},
		Spec: routev1.RouteSpec{
			Host: "apicast-wildcard.${WILDCARD_DOMAIN}",
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: "apicast-wildcard-router",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
			WildcardPolicy: routev1.WildcardPolicyType("${WILDCARD_POLICY}"),
			TLS: &routev1.TLSConfig{
				Termination:                   routev1.TLSTerminationType("edge"),
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyType("Allow")},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterService() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-wildcard-router",
			Labels: map[string]string{
				"app":                      "${APP_LABEL}",
				"3scale.component":         "apicast",
				"3scale.component-element": "wildcard-router",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       "http",
					Protocol:   v1.Protocol("TCP"),
					Port:       8080,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{"deploymentConfig": "apicast-wildcard-router"},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterDeploymentConfig() *appsv1.DeploymentConfig {
	return &appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps.openshift.io/v1", Kind: "DeploymentConfig"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-wildcard-router",
			Labels: map[string]string{
				"app":                      "${APP_LABEL}",
				"3scale.component":         "apicast",
				"3scale.component-element": "wildcard-router",
			},
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas: 1,
			Selector: map[string]string{
				"deploymentConfig": "apicast-wildcard-router",
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingParams: &appsv1.RollingDeploymentStrategyParams{
					IntervalSeconds: &[]int64{1}[0],
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
					TimeoutSeconds:      &[]int64{1800}[0],
					UpdatePeriodSeconds: &[]int64{1}[0],
				},
				Type: appsv1.DeploymentStrategyTypeRolling,
			},
			Triggers: appsv1.DeploymentTriggerPolicies{
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerType("ConfigChange"),
				},
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerType("ImageChange"),
					ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							"apicast-wildcard-router",
						},
						From: v1.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "amp-wildcard-router:latest",
						},
					},
				},
			},
			Template: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deploymentConfig":         "apicast-wildcard-router",
						"app":                      "${APP_LABEL}",
						"3scale.component":         "apicast",
						"3scale.component-element": "wildcard-router",
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "amp",
					Containers: []v1.Container{
						v1.Container{
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 8080,
									Protocol:      v1.ProtocolTCP,
									Name:          "http",
								},
							},
							Env:             wr.buildWildcardRouterEnv(),
							Image:           "amp-wildcard-router:latest",
							ImagePullPolicy: v1.PullPolicy("IfNotPresent"),
							Name:            "apicast-wildcard-router",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("500m"),
									v1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("120m"),
									v1.ResourceMemory: resource.MustParse("32Mi"),
								},
							},
							LivenessProbe: &v1.Probe{
								Handler: v1.Handler{TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("http"),
								}},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterEnv() []v1.EnvVar {
	return []v1.EnvVar{
		createEnvvarFromSecret("API_HOST", "system-master-apicast", "BASE_URL"),
	}
}
package operator

import (
	"fmt"
	"strconv"

	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	v1 "k8s.io/api/core/v1"
)

func (o *OperatorApicastOptionsProvider) GetApicastOptions() (*component.ApicastOptions, error) {
	optProv := component.ApicastOptionsBuilder{}
	optProv.AppLabel(*o.APIManagerSpec.AppLabel)
	optProv.TenantName(*o.APIManagerSpec.TenantName)
	optProv.WildcardDomain(o.APIManagerSpec.WildcardDomain)
	optProv.ManagementAPI(*o.ApicastSpec.ApicastManagementAPI)
	optProv.OpenSSLVerify(strconv.FormatBool(*o.ApicastSpec.OpenSSLVerify))        // TODO is this a good place to make the conversion?
	optProv.ResponseCodes(strconv.FormatBool(*o.ApicastSpec.IncludeResponseCodes)) // TODO is this a good place to make the conversion?

	o.setResourceRequirementsOptions(&optProv)
	o.setReplicas(&optProv)
	o.setNamespace(&optProv)
	o.setEnvironment(&optProv)
	res, err := optProv.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to create Apicast Options - %s", err)
	}
	return res, nil
}

func (o *OperatorApicastOptionsProvider) setResourceRequirementsOptions(b *component.ApicastOptionsBuilder) {
	if !*o.APIManagerSpec.ResourceRequirementsEnabled {
		b.ResourceRequirements(v1.ResourceRequirements{})
		// b.StagingResourceRequirements(v1.ResourceRequirements{})
		// b.ProductionResourceRequirements(v1.ResourceRequirements{})
	}
}

func (o *OperatorApicastOptionsProvider) setNamespace(b *component.ApicastOptionsBuilder) {
	if o.ApicastSpec.Namespace != nil {
		b.Namespace(*o.ApicastSpec.Namespace)
	} else {
		b.Namespace(o.Namespace)
	}
	// if o.APIManagerSpec.Apicast.ProductionSpec.Namespace == "" {
	// 	b.ProductionNamespace(o.Namespace)
	// } else {
	// 	b.ProductionNamespace(o.APIManagerSpec.Apicast.ProductionSpec.Namespace)
	// }

	// if o.APIManagerSpec.Apicast.StagingSpec.Namespace == "" {
	// 	b.StagingNamespace(o.Namespace)
	// } else {
	// 	b.StagingNamespace(o.APIManagerSpec.Apicast.StagingSpec.Namespace)
	// }
}

func (o *OperatorApicastOptionsProvider) setReplicas(b *component.ApicastOptionsBuilder) {
	b.Replicas(int32(*o.ApicastSpec.Replicas))
	// b.StagingReplicas(int32(*o.APIManagerSpec.Apicast.StagingSpec.Replicas))
	// b.ProductionReplicas(int32(*o.APIManagerSpec.Apicast.ProductionSpec.Replicas))
}

func (o *OperatorApicastOptionsProvider) setEnvironment(b *component.ApicastOptionsBuilder) {
	if o.ApicastSpec.Environment != nil {
		b.Environment(*o.ApicastSpec.Environment)
	}
	// b.StagingReplicas(int32(*o.APIManagerSpec.Apicast.StagingSpec.Replicas))
	// b.ProductionReplicas(int32(*o.APIManagerSpec.Apicast.ProductionSpec.Replicas))
}

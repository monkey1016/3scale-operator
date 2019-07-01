package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	"github.com/3scale/3scale-operator/pkg/common"
	templatev1 "github.com/openshift/api/template/v1"
)

type ImagesAdapter struct {
}

func NewImagesAdapter(options []string) Adapter {
	return NewAppenderAdapter(&ImagesAdapter{})
}

func (i *ImagesAdapter) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		templatev1.Parameter{
			Name:     "AMP_BACKEND_IMAGE",
			Required: true,
			Value:    "quay.io/3scale/3scale26:apisonator-3scale-2.6.0-ER1",
		},
		templatev1.Parameter{
			Name:     "AMP_ZYNC_IMAGE",
			Value:    "quay.io/3scale/3scale26:zync-3scale-2.6.0-ER1",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "AMP_APICAST_IMAGE",
			Value:    "quay.io/3scale/3scale26:apicast-3scale-2.6.0-ER1",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "AMP_SYSTEM_IMAGE",
			Value:    "quay.io/3scale/3scale26:porta-3scale-2.6.0-ER1",
			Required: true,
		},
		templatev1.Parameter{
			Name:        "ZYNC_DATABASE_IMAGE",
			Description: "Zync's PostgreSQL image to use",
			Value:       "centos/postgresql-10-centos7",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MEMCACHED_IMAGE",
			Description: "Memcached image to use",
			Value:       "memcached:1.5",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "IMAGESTREAM_TAG_IMPORT_INSECURE",
			Description: "Set to true if the server may bypass certificate verification or connect directly over HTTP during image import.",
			Value:       "false",
			Required:    true,
		},
	}
}

func (i *ImagesAdapter) Objects() ([]common.KubernetesObject, error) {
	imagesOptions, err := i.options()
	if err != nil {
		return nil, err
	}
	imagesComponent := component.NewAmpImages(imagesOptions)
	return imagesComponent.Objects(), nil
}

func (i *ImagesAdapter) options() (*component.AmpImagesOptions, error) {
	aob := component.AmpImagesOptionsBuilder{}
	aob.AppLabel("${APP_LABEL}")
	aob.AMPRelease("${AMP_RELEASE}")
	aob.ApicastImage("${AMP_APICAST_IMAGE}")
	aob.BackendImage("${AMP_BACKEND_IMAGE}")
	aob.SystemImage("${AMP_SYSTEM_IMAGE}")
	aob.ZyncImage("${AMP_ZYNC_IMAGE}")
	aob.ZyncDatabasePostgreSQLImage("${ZYNC_DATABASE_IMAGE}")
	aob.BackendRedisImage("${REDIS_IMAGE}")
	aob.SystemRedisImage("${REDIS_IMAGE}")
	aob.SystemMemcachedImage("${MEMCACHED_IMAGE}")

	aob.InsecureImportPolicy(false)

	return aob.Build()
}
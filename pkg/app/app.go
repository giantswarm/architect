package app

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	applicationv1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Namespace is the namespace where App CRs are created.
	namespace = "giantswarm"

	labelAppOperatorVersion = "app-operator.giantswarm.io/version"
	labelReleaseCyclePhase  = "release-operator.giantswarm.io/release-cycle-phase"
	labelServiceType        = "giantswarm.io/service-type"

	appOperatorVersion = "1.0.0"
)

// NewCR returns new application CR.
//
// AppCatalog is the name of the app catalog where the app stored.
func NewCR(name, appName, appVersion, appCatalog string) *applicationv1alpha1.App {
	appCR := &applicationv1alpha1.App{
		TypeMeta: applicationv1alpha1.NewAppTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				labelAppOperatorVersion: appOperatorVersion,
			},
		},
		Spec: applicationv1alpha1.AppSpec{
			Catalog: appCatalog,
			KubeConfig: applicationv1alpha1.AppSpecKubeConfig{
				InCluster: true,
			},
			Name:      appName,
			Namespace: namespace,
			Version:   appVersion,
		},
	}

	return appCR
}

func Print(w io.Writer, format string, appCR *applicationv1alpha1.App) error {
	var output []byte
	var err error

	switch format {
	case "json":
		output, err = json.Marshal(appCR)
		if err != nil {
			return microerror.Mask(err)
		}
	case "yaml":
		output, err = yaml.Marshal(appCR)
		if err != nil {
			return microerror.Mask(err)
		}
	default:
		return microerror.Maskf(wrongFormatError, "format: %q", format)
	}

	_, err = fmt.Fprintf(w, "%s", output)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

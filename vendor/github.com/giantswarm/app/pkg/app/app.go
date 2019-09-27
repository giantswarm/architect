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

// NewCR returns new application CR.
//
// AppCatalog is the name of the app catalog where the app stored.
func NewCR(name, appName, appVersion, appCatalog string) *applicationv1alpha1.App {
	appCR := &applicationv1alpha1.App{
		TypeMeta: applicationv1alpha1.NewAppTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "giantswarm",
			Labels: map[string]string{
				"app-operator.giantswarm.io/version": "1.0.0",
			},
		},
		Spec: applicationv1alpha1.AppSpec{
			Catalog: appCatalog,
			KubeConfig: applicationv1alpha1.AppSpecKubeConfig{
				InCluster: true,
			},
			Name:      appName,
			Namespace: "giantswarm",
			Version:   appVersion,
		},
	}

	return appCR
}

func Marshal(appCR *applicationv1alpha1.App, format string) (string, error) {
	var output []byte
	var err error

	switch format {
	case "json":
		output, err = json.Marshal(appCR)
		if err != nil {
			return "", microerror.Mask(err)
		}
	case "yaml":
		output, err = yaml.Marshal(appCR)
		if err != nil {
			return "", microerror.Mask(err)
		}
	default:
		return "", microerror.Maskf(executionFailedError, "format: %q", format)
	}

	return string(output), nil
}

func Print(w io.Writer, format string, appCR *applicationv1alpha1.App) error {
	output, err := Marshal(appCR, format)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = fmt.Fprintf(w, "%s", output)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

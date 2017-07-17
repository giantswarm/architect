package template

import (
	"fmt"
	"reflect"
	"strings"

	microerror "github.com/giantswarm/microkit/error"
)

const (
	// TemplateFieldFormat is the format for fields that SafeTemplate will template.
	TemplateFieldFormat = "{{ .%s }}"

	// StringTypeName is the name of the string type.
	StringTypeName = "string"
)

// safeTemplate takes a struct and a template byte array, and attempts to
// safely template the byte array.
// 'Safely' is defined as only templating fields of the format:
// `{{ .X }}`, where X is the name of some field in the supplied struct,
// and not smashing any other templating strings.
func SafeTemplate(s interface{}, template []byte) ([]byte, error) {
	if s == nil {
		return nil, microerror.MaskAny(nilTemplateStructError)
	}

	stringsToReplace := map[string]string{}

	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		value := v.FieldByName(name).String()

		typeName := v.Type().Field(i).Type.Name()
		if typeName != StringTypeName {
			return nil, microerror.MaskAnyf(notStringTypeError, "name: %v, type: %v", name, typeName)
		}

		stringsToReplace[name] = value
	}

	templateString := string(template)
	for name, value := range stringsToReplace {
		stringToReplace := fmt.Sprintf(TemplateFieldFormat, name)
		templateString = strings.Replace(templateString, stringToReplace, value, -1)
	}

	return []byte(templateString), nil
}

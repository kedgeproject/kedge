package pkg

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/runtime"

	"fmt"
	"os"

	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type App struct {
	Name           string `yaml:"name"`
	api_v1.PodSpec `yaml:",inline"`
}

func ReadFile(f string) ([]byte, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrap(err, "file reading failed")
	}
	return data, nil
}

func Convert(v *viper.Viper, cmd *cobra.Command) error {

	for _, file := range v.GetStringSlice("files") {
		d, err := ReadFile(file)
		if err != nil {
			return errors.New(err.Error())
		}

		var app App
		err = yaml.Unmarshal(d, &app)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal into internal struct")
		}

		runtimeObjects, err := CreateK8sObjects(&app)
		//fmt.Println("%s", spew.Sprint(runtimeObjects))

		for _, runtimeObject := range runtimeObjects {
			gvk, isUnversioned, err := api.Scheme.ObjectKind(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "ConvertToVersion failed")
			}
			if isUnversioned {
				return errors.New(fmt.Sprintf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject))
			}

			runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)

			data, err := yaml.Marshal(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			writeObject := func(o runtime.Object, data []byte) error {
				_, err := fmt.Fprintln(os.Stdout, "---")
				if err != nil {
					return errors.Wrap(err, "could not print to STDOUT")
				}

				_, err = os.Stdout.Write(data)
				return errors.Wrap(err, "could not write to STDOUT")
			}

			err = writeObject(runtimeObject, data)
			if err != nil {
				return errors.Wrap(err, "failed to write object")
			}
		}
	}

	return nil
}

func CreateK8sObjects(app *App) ([]runtime.Object, error) {

	var objects []runtime.Object

	p := api_v1.PodSpec(app.PodSpec)
	d := &ext_v1beta1.Deployment{
		ObjectMeta: api_v1.ObjectMeta{
			Name: app.Name,
			Labels: map[string]string{
				"app": app.Name,
			},
		},
		Spec: ext_v1beta1.DeploymentSpec{
			Template: api_v1.PodTemplateSpec{
				ObjectMeta: api_v1.ObjectMeta{
					Name: app.Name,
					Labels: map[string]string{
						"app": app.Name,
					},
				},
				Spec: p,
			},
		},
	}

	objects = append(objects, d)
	return objects, nil
}

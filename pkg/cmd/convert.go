package cmd

import (
	"io/ioutil"

	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/surajssd/kapp/pkg/encoding"
	"github.com/surajssd/kapp/pkg/transform/kubernetes"
	"k8s.io/client-go/pkg/runtime"
)

func Convert(files []string) error {

	for _, file := range files {

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "file reading failed")
		}

		app, err := encoding.Decode(data)
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal data")
		}

		ros, err := kubernetes.Transform(app)
		if err != nil {
			return errors.Wrap(err, "unable to convert data")
		}

		for _, runtimeObject := range ros {

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

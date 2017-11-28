package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

var Fixtures = os.ExpandEnv("$GOPATH/src/github.com/kedgeproject/kedge/tests/cmd/fixtures/")
var ProjectPath = "$GOPATH/src/github.com/kedgeproject/kedge/"
var BinaryLocation = os.ExpandEnv(ProjectPath + "kedge")
var imagename = "testrun"
var context = "kedge-build/"

func TestKedgeGenerate(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		input       string
		wantSuccess bool
		error       string
	}{
		{
			name:        "multiple configmaps given and name not specified",
			path:        "multi-configmapname/app.yaml",
			wantSuccess: false,
			error:       "unable to perform controller operations: unable to fix data: unable to fix ControllerFields: unable to fix configMaps: please specify name for app.configMaps[1]",
		},
		{
			name:        "generating deploymentconfig",
			path:        "controllers/input",
			input:       Fixtures + "controllers/os.yaml",
			wantSuccess: true,
		},
		{
			name:        "multiple containers given and name not specified for 2nd contaniner",
			path:        "multi-containername/nginx.yaml",
			wantSuccess: false,
			error:       "unable to perform controller operations: unable to fix data: unable to fix ControllerFields: unable to fix containers: please specify name for app.containers[1]",
		},
		{
			name:  "multiple volume claims with same name",
			path:  "multi-pvc-same-name/app.yaml",
			error: `unable to perform controller operations: unable to validate data: unable to validate controller fields: error validating volume claims: duplicate entry of volume claim "foo"`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			command := []string{"generate", "-f=" + Fixtures + tt.path}
			stdout, err := runCmd(t, command)
			if err != nil && !tt.wantSuccess {

				if stdout != tt.error {
					t.Fatalf("wanted error: %q\nbut got error: %q %v", tt.error, stdout, err)
				} else {
					t.Logf("failed with error: %q", stdout)
					return
				}
			} else if err != nil && tt.wantSuccess {
				t.Fatalf("wanted success, but test failed with error: %s %v", stdout, err)
			} else if err == nil && !tt.wantSuccess {
				t.Fatalf("expected to fail but passed")
			}

			// Read the data from the input file
			data, err := ioutil.ReadFile(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if diff := diff.Diff(string(data), fmt.Sprintf("%s\n", stdout)); diff != "" {
				t.Fatalf("wanted: \n%s\n======================\ngot: \n%s"+
					"\n======================\ndiff: %s", string(data), stdout, diff)
			}

		})
	}
}

func runCmd(t *testing.T, args []string) (string, error) {
	cmd := exec.Command(BinaryLocation, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	t.Logf("Running: %s", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		return strings.TrimSpace(stdout.String()), err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func Test_builderror(t *testing.T) {

	cmdStr := fmt.Sprintf("%s build -i %s -c %s", BinaryLocation, imagename, context)
	output, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		fmt.Println("Error executing command", err)
	}
	cmdStr = fmt.Sprintf("docker images | grep %s | awk '{print $1}'", imagename)
	output, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		fmt.Println("Error executing command", err)
	}
	if strings.TrimSpace(string(output)) != imagename {
		t.Errorf("Test Failed")
	}
}

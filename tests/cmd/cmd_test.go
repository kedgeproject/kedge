package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"path/filepath"

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
			error:       "unable to perform controller operations: unable to fix data: unable to fix configMaps: please specify name for app.configMaps[1]",
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
			error:       "unable to perform controller operations: unable to fix data: unable to fix deployments: unable to fix containers: please specify name for app.containers[1]",
		},
		{
			name:  "multiple volume claims with same name",
			path:  "multi-pvc-same-name/app.yaml",
			error: `unable to perform controller operations: unable to validate data: error validating volume claims: duplicate entry of volume claim "foo"`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			command := []string{"generate", "--skip-validation", "-f=" + Fixtures + tt.path}
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

func Test_stdin(t *testing.T) {

	kjson := `{"name": "httpd","deployments": [{"containers": [{"image": "centos/httpd"}]}]}`
	cmdStr := fmt.Sprintf("%s generate -f - <<EOF\n%s\nEOF\n", BinaryLocation, kjson)
	subproc := exec.Command("/bin/sh", "-c", cmdStr)
	output, err := subproc.Output()
	if err != nil {
		fmt.Println("Error executing command", err)
	}
	g, err := ioutil.ReadFile(Fixtures + "stdin/output.yml")
	if !bytes.Equal(output, g) {
		t.Errorf("Test Failed")
	}
}

func Test_examples(t *testing.T) {
	fileList := []string{}
	for _, dir := range []string{"examples", "docs/examples"} {
		err := filepath.Walk(os.ExpandEnv(ProjectPath)+dir, func(path string, f os.FileInfo, err error) error {
			_, file := filepath.Split(path)
			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
				if file != "cronjob.yaml" {
					fileList = append(fileList, path)
				}

			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	}

	for _, file := range fileList {
		cmdStr := fmt.Sprintf("%s generate -f %s", BinaryLocation, file)
		output, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
		if err != nil {
			t.Errorf("kedge generate failed for - %s\n Error is - %s", file, output)
		}
	}
}

func Test_init(t *testing.T) {

	//cmdStr := fmt.Sprintf("%s --name httpd --image centos/httpd --ports 80", BinaryLocation)
	svcname := "httpd"
	image := "centos/httpd"
	ports := "80"
	outputFile := os.ExpandEnv(ProjectPath) + "tests/cmd/kedge.yml"

	testCases := []struct {
		name        string
		command     []string
		wantSuccess bool
		error       string
		input       string
	}{
		{
			name:        "kedge init",
			command:     []string{"init"},
			wantSuccess: false,
			error:       "--name and --image are mandatory flags, Please provide these flags",
		},
		{
			name:        "kedge init with name",
			command:     []string{"init", "--name", svcname},
			wantSuccess: false,
			error:       "--name and --image are mandatory flags, Please provide these flags",
		},
		{
			name:        "kedge init with image",
			command:     []string{"init", "--image", image},
			wantSuccess: false,
			error:       "--name and --image are mandatory flags, Please provide these flags",
		},
		{
			name:        "kedge init with ports",
			command:     []string{"init", "--ports", ports},
			wantSuccess: false,
			error:       "--name and --image are mandatory flags, Please provide these flags",
		},
		{
			name:        "kedge init with name & image",
			command:     []string{"init", "--name", svcname, "--image", image},
			wantSuccess: true,
			input:       Fixtures + "init/kedge1.yml",
		},
		{
			name:        "kedge init with name & image & ports",
			command:     []string{"init", "--name", svcname, "--image", image, "--ports", ports},
			wantSuccess: true,
			input:       Fixtures + "init/kedge2.yml",
		},
		{
			name:        "kedge init with name & image & ports & restartPolicy",
			command:     []string{"init", "--name", svcname, "--image", image, "--ports", ports, "--restart-policy", "Always"},
			wantSuccess: true,
			input:       Fixtures + "init/kedge3.yml",
		},
		{
			name:        "kedge init with name & image & ports & restartPolicy & imagePullPolicy",
			command:     []string{"init", "--name", svcname, "--image", image, "--ports", ports, "--restart-policy", "Always", "--image-pull-policy", "Always"},
			wantSuccess: true,
			input:       Fixtures + "init/kedge4.yml",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			stdout, err := runCmd(t, tt.command)

			if err != nil && !tt.wantSuccess {
				if stdout != tt.error {
					t.Fatalf("Expected Error %s, But got %s", tt.error, stdout)
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
				t.Fatal(err, data)
			}

			outputData, err := ioutil.ReadFile(outputFile)
			if err != nil {
				t.Log("File not found")
			}
			if diff := diff.Diff(string(data), string(outputData)); diff != "" {
				t.Fatalf("wanted: \n%s\n======================\ngot: \n%s"+
					"\n======================\ndiff: %s", string(data), string(outputData), diff)
			}

			if _, err := os.Stat(outputFile); err == nil {
				err = os.Remove(outputFile)
				if err != nil {
					t.Log("Error in removing file")
				}
			}
		})
	}

}

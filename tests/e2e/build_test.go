package e2e

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

var err error

// Function to run shell commands
func runCmd(cmdS string) ([]byte, error) {
	var cmd *exec.Cmd
	var out, stdErr bytes.Buffer
	cmd = exec.Command("/bin/sh", "-c", cmdS)

	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running command %v: %s", cmd, err)

	}
	return out.Bytes(), nil
}

func DeleteRegistry(t *testing.T) {

	_, err = runCmd("docker stop -t 1 kedge_registry")
	if err != nil {
		t.Error(err)
	}

	_, err = runCmd("docker rm -v kedge_registry")
	if err != nil {
		t.Error(err)
	}

	t.Log("The registry has been deleted successfully")
}

// This creates a local docker registry,  runs `kedge build`, \
// and verifies if the image has been uploaded
func TestBuild(t *testing.T) {
	// Context Dir to be passed for kedge build
	var contextDir = ProjectPath + "/docs/examples/build"

	// Cleanup
	_, _ = runCmd("docker stop -t 1 kedge_registry")
	_, _ = runCmd("docker rm -v kedge_registry")

	endpoint := "unix:///var/run/docker.sock"
	cli, err := docker.NewClient(endpoint)
	if err != nil {
		t.Error(err)
	}

	dockerRegStart := "docker run -d -p 5000:5000 --restart=always --name kedge_registry registry:2"
	reg, err := runCmd(dockerRegStart)
	if err != nil {
		panic(err)
	}

	t.Logf("Local docker registry has been set up: %s", reg)

	defer DeleteRegistry(t)

	runKedgeCmd := BinaryLocation + " build -c " + contextDir + " -i localhost:5000/test:2.0 -p"

	t.Logf("Running '%s'\n", runKedgeCmd)
	kedgeBuild, err := runCmd(runKedgeCmd)
	if err != nil {
		t.Fatalf("Error running kedge: %v", err)
	}

	t.Log(string(kedgeBuild))

	// API call to list images
	images, err := cli.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		t.Error(err)
	}

	for _, image := range images {
		if image.RepoTags[0] == "localhost:5000/test:2.0" {
			t.Log("The image has been uploaded successfully.")

			rmc := "docker rmi " + image.ID

			_, err := runCmd(rmc)
			if err != nil {
				t.Error(err)
			}
		}
	}

}

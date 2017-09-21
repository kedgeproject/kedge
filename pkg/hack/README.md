This folder uses a "hack" of types.go from both OpenShift as well as Kubernetes
in order to prevent the full importation / vendoring of files from kubernetes/kubernetes
specifically: `github.com/kubernetes/kubernetes/blob/master/pkg/apis/apps/types.go` which
causes a conflict of vendored files between client-go and apimachinery.

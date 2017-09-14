# Development Guide

## Building Kedge

Read about building kedge [here](https://github.com/kedgeproject/kedge#building).

## Workflow
### Fork the main repository

1. Go to https://github.com/kedgeproject/kedge
2. Click the "Fork" button (at the top right)

### Clone your fork

The commands below require that you have $GOPATH. We highly recommended you put Kedge' code into your $GOPATH.

```console
git clone https://github.com/$YOUR_GITHUB_USERNAME/kedge.git $GOPATH/src/github.com/kedgeproject/kedge
cd $GOPATH/src/github.com/kedgeproject/kedge
git remote add upstream 'https://github.com/kedgeproject/kedge'
```

### Create a branch and make changes

```console
git checkout -b myfeature
# Make your code changes
```

### Keeping your development fork in sync

```console
git fetch upstream
git rebase upstream/master
```

Note: If you have write access to the main repository at github.com/kedgeproject/kedge, you should modify your git configuration so that you can't accidentally push to upstream:

```console
git remote set-url --push upstream no_push
```

### Committing changes to your fork

```console
git commit
git push -f origin myfeature
```

### Creating a pull request

1. Visit https://github.com/$YOUR_GITHUB_USERNAME/kedge.git
2. Click the "Compare and pull request" button next to your "myfeature" branch.
3. Check out the pull request process for more details

## `glide`, `glide-vc` and dependency management

Kedge uses `glide` to manage dependencies and `glide-vc` to clean vendor directory.
They are not strictly required for building Kedge but they are required when managing dependencies under the `vendor/` directory.
If you want to make changes to dependencies please make sure that `glide` and `glide-vc` are installed and are in your `$PATH`.

### Installing glide

There are many ways to build and host golang binaries. Here is an easy way to get utilities like `glide` and `glide-vc` installed:

Ensure that Mercurial and Git are installed on your system. (some of the dependencies use the mercurial source control system).
Use `apt-get install mercurial git` or `yum install mercurial git` on Linux, or `brew.sh` on OS X, or download them directly.

```console
go get -u github.com/Masterminds/glide
go get github.com/sgotti/glide-vc
```

Check that `glide` and `glide-vc` commands are working.

```console
glide --version
glide-vc -h
```

### Using glide

#### Adding new dependency

1. Update `glide.yaml` file

  Add new packages or subpackages to `glide.yaml` depending if you added whole
  new package as dependency or just new subpackage.

2. Get new dependencies

```bash
glide update --strip-vendor
```

3. Delete all unnecessary files from vendor

```bash
glide-vc --only-code --no-tests --use-lock-file
```

3. Commit updated glide files and vendor

```bash
git add glide.yaml glide.lock vendor
git commit
```


#### Updating dependencies

1. Set new package version in  `glide.yaml` file.

2. Clear cache

```bash
glide cc
```
This step is necessary if not done glide will pick up old data from it's cache.

3. Get new and updated dependencies

```bash
glide update --strip-vendor
```

4. Delete all unnecessary files from vendor

```bash
glide-vc --only-code --no-tests --use-lock-file
```

5. Commit updated glide files and vendor

```bash
git add glide.yaml glide.lock vendor
git commit
```

### PR review guidelines

- To merge a PR at least two LGTMs are needed to merge it

- If a PR is opened for more than two weeks, find why it is open for so long
if it is blocked on some other issue/pr label it as blocked and then also link
the issue it is blocked on. If it is outstanding for review and there are no
reviews on it ping maintainers.

- For PRs that have more than 500 LOC break it into pieces and merge it one
by one incrementally so that it is easy to review and going back and forth on
it is easier.

**Note**: Above guidelines are not hard rules use those with discretion

### Running tests

#### Run all tests except end-to-end tests

```bash
make test
```

#### Run end-to-end tests

Before running end to end tests locally make sure [minikube](https://github.com/kubernetes/minikube/)
is running.

```bash
make bin
make test-e2e
```

**Note**: When you run end to end tests, those tests are run in parallel. If
you are low on resources you can limit number of tests that run in parallel by
doing following:

```bash
make test-e2e PARALLEL=4
```

This will run only 4 tests in parallel. By default, it is set to the value of
`GOMAXPROCS`.

You may also add a timeout which will increase the overall timeout period for the tests.

```bash
make test-e2e TIMEOUT=15m
```

### spec.go conventions

- Add explanation on top of each struct and struct field in `spec.go` to explain what it does,
so that when OpenAPI spec is auto-generated it will show up there.

- Structs that are referred in any other struct in the form of an array should have a comment
of the format `// kedgeSpec: io.kedge.*`, where `*` is name of that struct. This becomes the
identity or reference of that struct in OpenAPI specification.

- If you are embedding a struct, there is no need to add an explanatory comment.

- For any struct that is embedded please add a k8s tag comment:
`// k8s: io.k8s.kubernetes.pkg.api.v1.ServicePort`.

- For all the fields that are optional please include a comment:
`// +optional`.

- Any struct that is defined in same file and is used in another struct, while embedding
please add a ref tag comment:
`// ref: io.kedge.ContainerSpec`.

- To find out what is the key or k8s reference tag for a particular struct in Kubernetes,
please refer to the swagger specification of Kubernetes for any particular release. For e.g
In Kubernetes 1.7, the reference tag for deployment is
`io.k8s.kubernetes.pkg.apis.apps.v1beta1.DeploymentSpec`.

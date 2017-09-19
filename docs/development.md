---
layout: default
permalink: /development/
redirect_from: "/docs/architecture.md"
---

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

### Validation

In order to facilitate consistent error messages, we ask that validation logic
adheres to the following guidelines whenever possible (though exceptional cases will exist).

* Be as precise as possible.
* Telling users what they CAN do is more useful than telling them what they
CANNOT do.
* When asserting a requirement in the positive, use "must".  Examples: "must be
greater than 0", "must match regex '[a-z]+'".  Words like "should" imply that
the assertion is optional, and must be avoided.
* When asserting a formatting requirement in the negative, use "must not".
Example: "must not contain '..'".  Words like "should not" imply that the
assertion is optional, and must be avoided.
* When asserting a behavioral requirement in the negative, use "may not".
Examples: "may not be specified when otherField is empty", "only `name` may be
specified".
* When referencing a literal string value, indicate the literal in
single-quotes. Example: "must not contain '..'".
* When referencing another field name, indicate the name in back-quotes.
Example: "must be greater than `request`".
* When specifying inequalities, use words rather than symbols.  Examples: "must
be less than 256", "must be greater than or equal to 0".  Do not use words
like "larger than", "bigger than", "more than", "higher than", etc.
* When specifying numeric ranges, use inclusive ranges when possible.

Taken from: [github.com/kubernetes/community/contributors/devel/api-conventions.md](https://github.com/kubernetes/community/blob/2bfe095e4dcd02b4ccd3e21c1f30591ca57518a6/contributors/devel/api-conventions.md#validation)


### Naming conventions

* Go field names must be CamelCase. JSON field names must be camelCase. Other
than capitalization of the initial letter, the two should almost always match.
No underscores nor dashes in either.
* Field and resource names should be declarative, not imperative (DoSomething,
SomethingDoer, DoneBy, DoneAt).
* Use `Node` where referring to
the node resource in the context of the cluster. Use `Host` where referring to
properties of the individual physical/virtual system, such as `hostname`,
`hostPath`, `hostNetwork`, etc.
* `FooController` is a deprecated kind naming convention. Name the kind after
the thing being controlled instead (e.g., `Job` rather than `JobController`).
* The name of a field that specifies the time at which `something` occurs should
be called `somethingTime`. Do not use `stamp` (e.g., `creationTimestamp`).
* We use the `fooSeconds` convention for durations, as discussed in the [units
subsection](#units).
  * `fooPeriodSeconds` is preferred for periodic intervals and other waiting
periods (e.g., over `fooIntervalSeconds`).
  * `fooTimeoutSeconds` is preferred for inactivity/unresponsiveness deadlines.
  * `fooDeadlineSeconds` is preferred for activity completion deadlines.
* Do not use abbreviations in the API, except where they are extremely commonly
used, such as "id", "args", or "stdin".
* Acronyms should similarly only be used when extremely commonly known. All
letters in the acronym should have the same case, using the appropriate case for
the situation. For example, at the beginning of a field name, the acronym should
be all lowercase, such as "httpGet". Where used as a constant, all letters
should be uppercase, such as "TCP" or "UDP".
* The name of a field referring to another resource of kind `Foo` by name should
be called `fooName`. The name of a field referring to another resource of kind
`Foo` by ObjectReference (or subset thereof) should be called `fooRef`.
* More generally, include the units and/or type in the field name if they could
be ambiguous and they are not specified by the value or value type.
* The name of a field expressing a boolean property called 'fooable' should be
called `Fooable`, not `IsFooable`.

Taken from: [github.com/kubernetes/community/contributors/devel/api-conventions.md](https://github.com/kubernetes/community/blob/2bfe095e4dcd02b4ccd3e21c1f30591ca57518a6/contributors/devel/api-conventions.md#naming-conventions)

### Optional vs. Required

Fields must be either optional or required.

Optional fields have the following properties:

- They have the `+optional` comment tag in Go.
- They are a pointer type in the Go definition (e.g. `bool *awesomeFlag`) or
have a built-in `nil` value (e.g. maps and slices).

In most cases, optional fields should also have the `omitempty` struct tag (the 
`omitempty` option specifies that the field should be omitted from the json
encoding if the field has an empty value).


Required fields have the opposite properties, namely:

- They do not have an `+optional` comment tag.
- They do not have an `omitempty` struct tag.
- They are not a pointer type in the Go definition (e.g. `bool otherFlag`).

Using the `+optional` or the `omitempty` tag causes OpenAPI documentation to 
reflect that the field is optional.

Using a pointer allows distinguishing unset from the zero value for that type.
There are examples of this in the codebase. However:

- it can be difficult for implementors to anticipate all cases where an empty
value might need to be distinguished from a zero value
- having a pointer consistently imply optional is clearer

Therefore, we ask that pointers always be used with optional fields that do not
have a built-in `nil` value.

Inspired from: [github.com/kubernetes/community/contributors/devel/api-conventions.md](https://github.com/kubernetes/community/blob/2bfe095e4dcd02b4ccd3e21c1f30591ca57518a6/contributors/devel/api-conventions.md#optional-vs-required)

### General guidelines for developers

- When you add a new function/method

  - Add unit-tests

- When you add a new feature

  - Add an example in docs/example with its explanation README, for e.g. [health](https://github.com/kedgeproject/kedge/tree/master/docs/examples/health).
  - Add an e2e test on above example, for e.g. see test code for [health](https://github.com/kedgeproject/kedge/blob/cfee15ffde02c611d08420699a43869706be2d53/tests/e2e/e2e_test.go#L272).
  - Add this feature information to file-reference, for e.g. see [health section](https://github.com/kedgeproject/kedge/blob/master/docs/file-reference.md#health).

### golang dependency import conventions

Imports MUST be arranged in three sections, separated by an empty line.

```go
stdlib
kedge
thirdparty
```

For example:

```go
"fmt"
"io"
"os/exec"

pkgcmd "github.com/kedgeproject/kedge/pkg/cmd"
"github.com/kedgeproject/kedge/pkg/spec"

"github.com/ghodss/yaml"
"github.com/pkg/errors"
```

Once arranged, let `gofmt` sort the sequence of imports.

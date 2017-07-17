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
1. Update `glide.yml` file.

  Add new packages or subpackages to `glide.yml` depending if you added whole new package as dependency or
  just new subpackage.

2. Run `glide update --strip-vendor` to get new dependencies.
   Than run `glide-vc --only-code --no-tests` to delete all unnecessary files from vendor.

3. Commit updated `glide.yml`, `glide.lock` and `vendor` to git.


#### Updating dependencies

1. Set new package version in  `glide.yml` file.

2. Run `glide update --strip-vendor` to update dependencies.
   Than run `glide-vc --only-code --no-tests` to delete all unnecessary files from vendor.


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

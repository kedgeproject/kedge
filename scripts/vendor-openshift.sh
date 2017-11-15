#!/bin/bash 

# 'glide update -v' has to be run before running this script.
# If you do 'glide update' after running this scripts, glide is going to ovewrite all changes.


# This adds OpenShift and all packages that are vendored by OpenShift
# to project vendor directory. It has to be run from root directory of the project
# (where vendor directory is).

# For more information why we are doing this see comments glide.yaml


OPENSHIFT_REPO="https://github.com/openshift/origin"
OPENSHIFT_VERSION="v3.6.1"

PROJECT_VENDOR="./vendor"

TMP_OPENSHIFT=`mktemp -d`

echo "Cloning OpenShift $OPENSHIFT_VERSION"
git clone --branch $OPENSHIFT_VERSION --depth 1 $OPENSHIFT_REPO $TMP_OPENSHIFT


# How deep is the package in directory structure?
# example: package hosted in github.com has three level strucure (github.com/<namespace>/<pkgname>)
#          packages in k8s.io has only two level strucutre (k8s.io/<pkgname>)
# we need to know where it is so we can move whole package (nothing more nothing less)
#
# If we were to move whole github.com than whole github.com gets replaced in target directory,
# and if target directory had some extra libs from gihtub.com that are not in openshift vendor
# they would get removed.

TWO_LEVEL="cloud.google.com go.pedge.io go4.org google.golang.org gopkg.in k8s.io vbom.ml"
THREE_LEVEL="bitbucket.org github.com golang.org"

function movePKG {
    # This function moves package from OpenShift vendor to this project vendor 
    # Takes one positional argument - package name (github.com/foo/bar)
    pkg=$1

    target="$PROJECT_VENDOR/$pkg"
    target_parent=`dirname $target`

    rm -rf $target
    mkdir -p $target_parent

    echo "Moving $pkg from OpenShift vendor to this project vendor directory"
    mv -f $TMP_OPENSHIFT/vendor/$pkg $target_parent    
}

# check if we cover everything from openshift vendor
# every domain in OpenShift vendor has to be covered in lists above (to define if its 2 or 3 level structure)
for path in `find $TMP_OPENSHIFT/vendor -maxdepth 1 -mindepth 1 -type d | sort`; do
    domain=`basename $path`

    found=false
    for t in $TWO_LEVEL $THREE_LEVEL; do
        if [ "$t" == "$domain" ]; then
            found=true
        fi
    done

    if [ $found == false ]; then
        echo "ERROR: structure for $domain is not defined"
        exit 1
    fi

done


# move packages from OpenShifts vendor dir to project vendor dir

# for every package that is organized in two level directory strucure
for domain in $TWO_LEVEL; do
    # on Linux you can use just find with  `-printf "$domain/%P\n"` instead of awk, but this is not available on MacOS
    pkgs=`find -L "${TMP_OPENSHIFT}/vendor/$domain"  -maxdepth 1 -mindepth 1 -type d | awk -F/ '{ print($(NF-1)"/"$(NF)) }'`
    for pkg in $pkgs; do
        movePKG $pkg
    done
done

# for every package that is organized in three level directory strucure
for domain in $THREE_LEVEL; do
    # on Linux you can use just find with  `-printf "$domain/%P\n"` instead of awk, but this is not available on MacOS
    pkgs=`find -L "${TMP_OPENSHIFT}/vendor/$domain"  -maxdepth 2 -mindepth 2 -type d | awk -F/ '{ print($(NF-2)"/"$(NF-1)"/"$(NF)) }'`
    for pkg in $pkgs; do
        movePKG $pkg
    done
done


# OpenShift vendor directory shouldn't contain any *.go files, they should be all moved to project vendor directory
remaining_go_files=`find ${TMP_OPENSHIFT}/vendor -type f -name *.go`
if [[ "`$remaining_go_files | wc -l`" -ne 0 ]]; then
    echo `echo $remaining_go_files| wc -l`
    echo "ERROR: There *.go files remaining in OpenShift vendor directory"
    echo $remaining_go_files
    exit 1
fi

# clean up openshift
rm -rf $TMP_OPENSHIFT/vendor $TMP_OPENSHIFT/.git

# now move OpenShift itself to project vendor
echo "Moving OpenShift to project vendor directory"
rm -rf $PROJECT_VENDOR/github.com/openshift/origin
mkdir -p $PROJECT_VENDOR/github.com/openshift/
mv $TMP_OPENSHIFT $PROJECT_VENDOR/github.com/openshift/origin

echo "DONE."



#!/bin/sh

DATE=`date --iso-8601=date`
TIME=`date --iso-8601=seconds`


# generate ./bin/info.txt
# this file contains information about what files were build and whan

commit_id=$(git rev-parse HEAD)
origin=$(git config --get remote.origin.url)

cat > "./bin/info.txt" <<EOF
date: ${TIME}
build_from: ${origin}
commit_id: ${commit_id}
files:
EOF

for f in $(ls -1 ./bin/* | grep -v info.txt); do
  sha256sum=$(sha256sum $f | cut -d ' ' -f 1);
  name=$(basename $f)
  updated_on=$(stat -c %y $f)
  cat >> "./bin/info.txt" <<EOF
  - name: ${name}"
    sha256sum: ${sha256sum}"
    updated_on: ${updated_on}
EOF
done



# generate .bintray.json
# this file contains all information on what will be upload to bintray
# for mor info: https://docs.travis-ci.com/user/deployment/bintray/
cat > "./.bintray.json" <<EOF
{
    "package": {
        "name": "kedge",
        "repo": "kedge",
        "subject": "kedgeproject",
        "desc": "Concise Application Definition for Kubernetes",
        "website_url": "https://github.com/kedgeproject/kedge",
        "issue_tracker_url": "https://github.com/kedgeproject/kedge/issues",
        "vcs_url": "https://github.com/kedgeproject/kedge.git",
        "licenses": ["Apache-2.0"],
        "public_download_numbers": false,
        "public_stats": false
    },

    "version": {
        "name": "latest",
        "desc": "Kedge build from master branch",
        "released": "${DATE}",
        "vcs_tag": "${TRAVIS_COMMIT}",
        "attributes": [{"name": "TRAVIS_JOB_NUMBER", "values" : ["${TRAVIS_JOB_NUMBER}"], "type": "string"},
                       {"name": "TRAVIS_JOB_ID", "values" : ["${TRAVIS_JOB_ID}"], "type": "string"},
                       {"name": "TRAVIS_COMMIT", "values" : ["${TRAVIS_COMMIT}"], "type": "string"},
                       {"name": "TRAVIS_BRANCH", "values" : ["${TRAVIS_BRANCH}"], "type": "string"},
                       {"name": "TRAVIS_PULL_REQUEST", "values" : ["${TRAVIS_PULL_REQUEST}"], "type": "string"},
                       {"name": "date", "values" : ["${TIME}"], "type": "date"}],
        "gpgSign": false
    },

    "files":
        [
            {"includePattern": "bin/(.*)",
             "uploadPattern": "./latest/\$1", 
             "matrixParams": {"override": 1 }
            }
        ],
    "publish": true
}
EOF
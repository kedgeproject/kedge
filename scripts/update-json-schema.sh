#!/usr/bin/env bash

# Ensures that we run on Travis
if [ "$TRAVIS_BRANCH" != "master" ] || [ "$GENERATE_JSON_SCHEMA" != "yes" ] || [ "$TRAVIS_SECURE_ENV_VARS" == "false" ] || [ "$TRAVIS_PULL_REQUEST" != "false" ] ; then
    echo "Must be: a merged PR on the master branch, GENERATE_JSON_SCHEMA=yes, TRAVIS_SECURE_ENV_VARS=false"
    exit 0
fi

KEDGE_REPO_NAME="kedge-json-schema"
GENERATOR_REPO="git@github.com:kedgeproject/json-schema.git"
DEPLOY_KEY="scripts/deploy_key"
KEDGE_JSON_SCHEMA_IMAGE="containscafeine/kedge-json-schema:latest"
GIT_USER="kedge-bot"
GIT_EMAIL="shubham@linux.com"
OUTPUT_DIR="master"

# decrypt the private key
openssl aes-256-cbc -K $encrypted_91569b511922_key -iv $encrypted_91569b511922_iv -in "$DEPLOY_KEY.enc" -out "$DEPLOY_KEY" -d
chmod 600 "$DEPLOY_KEY"
eval `ssh-agent -s`
ssh-add "$DEPLOY_KEY"

# clone the JSON Schema Generator repo
git clone "$GENERATOR_REPO" "$KEDGE_REPO_NAME"
mkdir -p "$KEDGE_REPO_NAME"/"$OUTPUT_DIR"
cd "$KEDGE_REPO_NAME"/"$OUTPUT_DIR"
docker run --rm -v `pwd`:/data:Z "$KEDGE_JSON_SCHEMA_IMAGE"

# add relevant user information
git config user.name "$GIT_USER"

# email assigned
git config user.email "$GIT_EMAIL"
git add --all

# Check if anything changed, and if it's the case, push to origin/master.
if git commit -m 'Update JSON Schema' -m "Commit: https://github.com/kedgeproject/kedge/commit/$TRAVIS_COMMIT" ; then
  git push
fi

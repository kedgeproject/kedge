#!/bin/bash

set -e

echo "Generating git config"
echo '
[user]
  name = Kedge Schema Bot
  email = kedgeschema@gmail.com
' | tee ~/.gitconfig

echo "Cloning Kedge JSONSchema repository"
cd
yes | git clone git@github.com:kedgeproject/json-schema.git
cd json-schema
docker run -v `pwd`:/data:Z surajd/kedgeschema

echo "Pushing all the generated content to github"
git add .
git commit -m "auto generated on $(date)"
git push origin master

echo "Schema pushed successfully"

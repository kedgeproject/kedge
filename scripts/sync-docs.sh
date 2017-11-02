#!/usr/bin/env bash

# Ensures that we run on Travis
if [ "$TRAVIS_BRANCH" != "master" ] || [ "$BUILD_DOCS" != "yes" ] || [ "$TRAVIS_SECURE_ENV_VARS" == "false" ] || [ "$TRAVIS_PULL_REQUEST" != "false" ] ; then
    echo "Must be: a merged pr on the master branch, BUILD_DOCS=yes, TRAVIS_SECURE_ENV_VARS=false"
    exit 0
fi

DOCS_REPO_NAME="kedge"
DOCS_REPO_URL="git@github.com:kedgeproject/kedge.git"
DOCS_KEY="scripts/deploy_key"
DOCS_USER="kedge-bot"
DOCS_EMAIL="cdrage+kedge@redhat.com"
DOCS_BRANCH="gh-pages"
DOCS_FOLDER="docs"

# decrypt the private key
openssl aes-256-cbc -K $encrypted_91569b511922_key -iv $encrypted_91569b511922_iv -in "$DOCS_KEY.enc" -out "$DOCS_KEY" -d
chmod 600 "$DOCS_KEY"
eval `ssh-agent -s`
ssh-add "$DOCS_KEY"

# clone the repo
git clone "$DOCS_REPO_URL" "$DOCS_REPO_NAME"

# change to that directory (to prevent accidental pushing to master, etc.)
cd "$DOCS_REPO_NAME"

# switch to gh-pages and grab the docs folder from master
git checkout gh-pages
git checkout master docs

# Remove README.md from docs folder as it isn't relevant
rm docs/README.md

# Use introduction.md instead as the main index page
mv docs/introduction.md index.md


# Check that index.md has the appropriate Jekyll format
index="index.md"
if cat $index | head -n 1 | grep "\-\-\-";
then
echo "index.md already contains Jekyll format"
else
# Remove ".md" from the name
name=${index::-3}
echo "Adding Jekyll file format to $index"
jekyll="---
layout: default
---
"
echo -e "$jekyll\n$(cat $index)" > $index
fi

# clean-up the docs and convert to jekyll-friendly docs
cd docs
for filename in *.md; do
    if cat $filename | head -n 1 | grep "\-\-\-";
    then
    echo "$filename already contains Jekyll format"
    else
    # Remove ".md" from the name
    name=${filename::-3}
    echo "Adding Jekyll file format to $filename"
    jekyll="---
layout: default
permalink: /$name/
redirect_from: 
  - /docs/$name.md/
---
"
    echo -e "$jekyll\n$(cat $filename)" > $filename
    fi
done
cd ..

# Build slate / file-reference
cp mv docs/file-reference.md slate/source/index.html.md
slate="---
title: Kedge File Reference

language_tabs:
  - yaml

toc_footers:
  - <a href='http://kedgeproject.org'>kedgeproject.org</a>
  - <a href='https://github.com/kedgeproject/kedge'>Kedge on GitHub</a>

search: true
---
"
echo -e "$slate"\n$(cat slate/source/index.html.md)" >  slate/source/index.html.md
cd slate
docker run --rm -v $PWD:/usr/src/app/source -w /usr/src/app/source cdrage/slate bundle exec middleman build --clean
rm _site/file-reference
mv build _site/file-reference
cd ..

# add relevant user information
git config user.name "$DOCS_USER"

# email assigned
git config user.email "$DOCS_EMAIL"
git add --all

# Check if anything changed, and if it's the case, push to origin/master.
if git commit -m 'Update docs' -m "Commit: https://github.com/kedgeproject/kedge/commit/$TRAVIS_COMMIT" ; then
  git push
fi

# cd back to the original root folder
cd ..


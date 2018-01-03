#!/usr/bin/env bash

curl -O https://raw.githubusercontent.com/kedgeproject/json-schema/master/master/kedge-json-schema.json
node scripts/modify-schema.js 
echo -e 'package validation\n\nvar SchemaJson = `' > kedgeschema.go
sed -i 's/`//g' schema.json
cat schema.json >> kedgeschema.go
sed -i -e '$a`' kedgeschema.go
sed -i -e '5s/$/   "additionalProperties": false,/' kedgeschema.go
mv kedgeschema.go pkg/validation/
rm kedge-json-schema.json schema.json


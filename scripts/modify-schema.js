var fs = require('fs')
var jsonSchemaObj = JSON.parse(fs.readFileSync("./kedge-json-schema.json").toString());

// In case of deploymentconfig
a = jsonSchemaObj.properties.deploymentConfigs.items.required
remove(a, ["test", "replicas", "stratergy", "triggers"])
jsonSchemaObj.properties.deploymentConfigs.items.required = a;

//In case of  buildconfig
a = jsonSchemaObj.properties.buildConfigs.items.required
remove(a, ["nodeSelector", "stratergy", "triggers"])
jsonSchemaObj.properties.buildConfigs.items.required = a;
//console.log(JSON.stringify(jsonSchemaObj, null,"\t"));

//In case of  imagestream
a = jsonSchemaObj.properties.imageStreams.items.properties.tags.items.required
remove(a, ["annotations", "generation"])
jsonSchemaObj.properties.imageStreams.items.properties.tags.items.required = a;

//In case of  routes - case1
a = jsonSchemaObj.properties.routes.items.required
remove(a, ["host"])
jsonSchemaObj.properties.routes.items.required = a;

//In case of  routes - case2
a = jsonSchemaObj.properties.routes.items.properties.to.required
remove(a, ["weight"])
jsonSchemaObj.properties.routes.items.properties.to.required = a;

p = JSON.stringify(jsonSchemaObj, null,"\t")

fs.writeFileSync("schema.json", p)
//console.log(JSON.stringify(jsonSchemaObj, null,"\t"));


function remove(array, elementArraY){
 for (let index = 0; index < elementArraY.length; index++) {
     array.splice(array.indexOf(elementArraY[index]), 1);
     
 }

}

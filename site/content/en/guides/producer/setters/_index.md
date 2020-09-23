---
title: "Create Setters"
linkTitle: "Setters"
weight: 2
type: docs
description: >
    Create high-level setters to provide imperative configuration editing
    commands.
---
Setters provide a solution for template-free setting or substitution of field
values through package metadata (OpenAPI).  They are a safer alternative to
other substitution techniques which do not have the context of the
structured data -- e.g. using `sed` to replace values.

The OpenAPI definitions for setters are defined in a Kptfile and referenced by
a fields through comments on the fields.

Setters may be invoked to programmatically modify the configuration
using `kpt cfg set` to set and/or substitute values.

{{% pageinfo color="primary" %}}
Creating a setter requires that the package has a Kptfile.  If one does
not exist for the package, run `kpt pkg init DIR/` to create one.
{{% /pageinfo %}}

## Setters explained

Following is a short explanation of the command that will be demonstrated
in this guide.

### Data model

- Fields reference setters through OpenAPI definitions specified as
  line comments -- e.g. `# { "$kpt-set": "replicas-setter" }`
- OpenAPI definitions are provided through the Kptfile

### Command control flow

1. Read the package Kptfile and resources.
2. Change the setter OpenAPI value in the Kptfile
3. Locate all fields which reference the setter and change their values.
4. Write both the modified Kptfile and resources back to the package.

{{< svg src="images/set-command" >}}

#### Creating a Setter

Setters may be created either manually (by editing the Kptfile directly), or
programmatically (through the `create-setter` command).  The `create-setter`
command will:

1. create a new OpenAPI definition for a setter in the Kptfile
2. create references to the setter definition on the resource fields

```yaml
# deployment.yaml -- original
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3
```

```sh
# create or update a setter named "replicas"
# match fields with the value "3"
kpt cfg create-setter hello-world/ replicas 3
```

```yaml
# Kptfile -- updated
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      x-k8s-cli:
        setter:
          name: "replicas"
          value: "3"
```

```yaml
# deployment.yaml -- updated
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

#### Targeting fields using the path

The basic way to create a setter only matches fields based on the value. But
in some situations there might be several fields with the same value, but not all
of them should be targeted by the setter. In these situations, use the `--field`
flag to provide either the name of the field, the full path to the field, or a
partial (suffix) path to the field. Only fields that match both the path and
the value will be targeted by the setter.

```sh
# create a setter named "replicas" and but only target the field name replicas in the spec
kpt cfg create-setter hello-world/ replicas 3 --field="spec.replicas"
```

```yaml
# deployment-foo.yaml
kind: Deployment
metadata:
  name: foo
  annotations:
    replicas: 3
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

The path is always just the names of the properties on the path to the field,
regardless whether the field is nested inside a sequence. Targeting specific
elements inside a sequence is not supported.

```sh
# create a setter named "mountName" and only target the name of the volume mount
kpt cfg create-setter hello-world/ mountName nginx --field="spec.containers.volumeMounts.name"
```

```yaml
# deployment-nginx.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: nginx # {"$kpt-set":"mountName"}
      mountPath: /usr/share/nginx
    - name: temp
      mountPath: /tmp
```

#### Auto setters

Some setters are automatically set deriving the values from the output of `gcloud config list`
command, when the package is fetched using `kpt pkg get`. Package consumers need not
invoke `kpt cfg set` on them explicitly. Following are the names of supported auto-setters:

```
gcloud.core.project
gcloud.project.projectNumber
gcloud.compute.region
gcloud.compute.zone
```

#### Invoking a Setter

```yaml
# deployment.yaml -- original
kind: Deployment
metadata:
 name: helloworld-gke
 labels:
   app: hello
spec:
 replicas: 3 # {"$kpt-set":"replicas"}
```

```sh
# set the replicas field to 5
kpt cfg set DIR/ replicas 5
```

```yaml
# deployment.yaml -- updated
kind: Deployment
metadata:
 name: helloworld-gke
 labels:
   app: hello
spec:
 replicas: 5 # {"$kpt-set":"replicas"}
```

#### OpenAPI Validations

Users can input any additional validation constraints during `create-setter`
operation in the form of openAPI schema. Relevant openAPI specification
constraints can be provided in json file format. The `set` operation validates
the input value against provided schema during setter creation and throws an
error if the input value doesn't meet any of the constraints.

```sh
$ cat /path/to/file.json
{"maximum": 10, "type": "integer"}

# create setter with openAPI property constraints
kpt cfg create-setter DIR/ replicas 3 --schema-path /path/to/file.json
```

The command creates setter with the following definition

```yaml
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      maximum: 10
      type: integer
      x-k8s-cli:
        setter:
          name: replicas
          value: "3"
```

```sh
# try to set value not adhering to the constraints
kpt cfg set DIR/ replicas 11
```

```sh
error: The input value doesn't validate against provided OpenAPI schema:
validation failure list: replicas in body should be less than or equal to 10
```

```sh
Example schema for integer validation

{
  "maximum": 10,
  "type": "integer",
  "minimum": 3,
  "format": "int64",
  "multipleOf": 2
}

Example schema for string validation

{
  "maxLength": 10,
  "type": "string",
  "minLength": 3,
  "pattern": "^[A-Za-z]+$",
  "enum": [
    "nginx",
    "ubuntu"
  ]
}

Example schema for array validation

{
  "maxItems": 10,
  "type": "array",
  "minItems": 3,
  "uniqueItems": true,
  "items": {
    "type": "string",
    "maxLength": 4
  }
}

```

Relevant resources for more information: [OpenAPI types]

#### Setting Lists

It is possible to create setters for fields which are a list of strings/integers.
The setter type must be `array`, and the reference must be on the list field.
The list setter will take variable args for its value rather than a single value.

**Note:** You should skip passing the value arg while creating array setters. `field`
flag is required for array setters.

```yaml
# example.yaml
apiVersion: example.com/v1beta1
kind: Example
spec:
  list:
  - "a"
  - "b"
```

```yaml
# Kptfile
kind: Kptfile
```

`$ kpt cfg create-setter DIR/ list --type array --field spec.list`

```yaml
# example.yaml
apiVersion: example.com/v1beta1
kind: Example
spec:
  list: # {"$kpt-set":"list"}
  - "a"
  - "b"
```

```yaml
# Kptfile
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.list:
      type: array
      x-k8s-cli:
        setter:
          name: list
          listValues:
          - "a"
          - "b"
```

`$ kpt cfg set DIR/ list c d e`

```yaml
# example.yaml
apiVersion: example.com/v1beta1
kind: Example
spec:
  list: # {"$kpt-set":"list"}
  - "c"
  - "d"
  - "e"
```

```yaml
# Kptfile
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.list:
      type: array
      x-k8s-cli:
        setter:
          name: list
          listValues:
          - "c"
          - "d"
          - "e"
```

#### Enumerations

Setters may be configured to map an enum input to a different value set
in the configuration.

e.g. users set `small`, `medium`, `large` cpu sizes, and these are mapped
to numeric values set in the configuration.

This may be done by modifying the Kptfile OpenAPI definitions as shown here:

```yaml
openAPI:
  definitions:
    io.k8s.cli.setters.cpu:
      x-k8s-cli:
        setter:
          name: "cpu"
          value: "small"
          # enumValues will replace the user provided key with the
          # map value when setting fields.
          enumValues:
            small: "0.5"
            medium: "2"
            large: "4"
```

Set would change the configuration like this:

```yaml
kind: Deployment
metadata:
  name: foo
spec:
  template:
    spec:
      containers:
      - name: foo
    resources:
      requests:
        cpu: "0.5" # {"$kpt-set":"cpu"}
```

#### Marking a field as required

Package publisher can mark a setter as required to convey the consumer that the
setter value must be set before triggering live apply/preview operation
on the package.

```sh
# create a setter named "replicas" and mark it as required
kpt cfg create-setter hello-world/ replicas 3 --required
```

```yaml
# deployment-foo.yaml
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

```yaml
# Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      x-k8s-cli:
        setter:
          name: "replicas"
          value: "3"
          required: true
```

```sh
# if you live apply/preview without setting the value
kpt live apply hello-world/

error: setter replicas is required but not set, please set it and try again
```

```sh
# set the replicas value
kpt cfg set hello-world/ replicas 4

kpt live apply hello-world/
# Success
```

#### Modifying a Setter

Setters are uniquely identified by their name. create-setter command can be leveraged
to modify an existing setter definition. Users may choose to modify setters in following
scenarios.

1. Add new features to existing setters, such as openAPI validations, marking setter
as required setter etc.
2. Add existing setter references to new resources that are added to package.

Consider an existing package with deployment-foo.yaml and Kptfile as follows.

```yaml
# deployment-foo.yaml
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

```yaml
# Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      x-k8s-cli:
        setter:
          name: "replicas"
          value: "3"
```

Add new resource to the package

```yaml
# deployment-bar.yaml
kind: Deployment
metadata:
  name: bar
spec:
  replicas: 3
```

```sh
$ cat /path/to/file.json
{"maximum": 10, "type": "integer"}

# modify existing setter replicas with openAPI property constraints and required flag
kpt cfg create-setter DIR/ replicas 3 --schema-path /path/to/file.json --required
```

```yaml
# Kptfile -- updated
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      maximum: 10
      type: integer
      x-k8s-cli:
        setter:
          name: replicas
          value: "3"
          required: true
```

```yaml
# deployment-foo.yaml -- remains unchanged
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

```yaml
# deployment-bar.yaml -- setter ref added
kind: Deployment
metadata:
  name: bar
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

#### Deleting a Setter

Setters may be deleted either manually (by editing the Kptfile directly), or
programmatically (through the `delete-setter` command).  The `delete-setter`
command will:

1. delete an OpenAPI definition for a setter in the Kptfile
2. remove references to the setter definition on the resource fields

```yaml
# Kptfile -- original
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      x-k8s-cli:
        setter:
          name: "replicas"
          value: "3"
```

```yaml
# deployment.yaml -- original
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3 # {"$kpt-set":"replicas"}
```

```sh
# delete a setter named "replicas"
kpt cfg delete-setter replicas
```

```yaml
# Kptfile -- updated
openAPI:
```

```yaml
# deployment.yaml -- updated
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 3
```

[OpenAPI types]: https://swagger.io/docs/specification/data-models/data-types/
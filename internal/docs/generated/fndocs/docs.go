// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "mdtogo"; DO NOT EDIT.
package fndocs

var READMEShort = `Generate, transform, or validate configuration files using containerized functions.`
var READMELong = `
Functions are executables ([that you can write](#developing-functions)) packaged in container images which accept a collection of
Resource configuration as input, and emit a collection of Resource configuration as output.

| Command  | Description                            |
|----------|----------------------------------------|
| [run]    | locally execute one or more functions  |
| [sink]   | explicitly specify an output sink      |
| [source] | explicitly specify an input source     |

**Data Flow**:  local configuration or stdin -> kpt [fn] (runs a container) -> local configuration or stdout

| Configuration Read From | Configuration Written To |
|-------------------------|--------------------------|
| local files or stdin    | local files or stdout    |

Functions may be used to:

* Generate configuration from templates, DSLs, CRD-style abstractions, key-value pairs, etc.-- e.g. expand Helm charts, JSonnet, Jinja, etc.
* Inject fields or otherwise modify configuration -- e.g.add init-containers, side-cars, etc
* Rollout configuration changes across an organization -- e.g.similar to
  https://github.com/reactjs/react-codemod
* Validate configuration -- e.g.ensure organizational policies are enforced

Functions may be run at different times depending on the function and the organizational needs:

* as part of the build and development process
* as pre-commit checks
* as PR checks
* as pre-release checks
* as pre-rollout checks

#### Functions Catalog

[KPT Functions Catalog][catalog] repository documents a catalog of kpt functions implemented using different toolchains.

#### Developing Functions


| Language   | Documentation               | Examples                    |
|------------|-----------------------------|-----------------------------|
| Typescript | [KPT Functions SDK][sdk-ts] | [examples][sdk-ts-examples] |
| Go         | [kustomize/kyaml][kyaml]    | [example][kyaml-example]    |
`
var READMEExamples = `
    # run the function defined by gcr.io/example.com/my-fn as a local container
    # against the configuration in DIR
    kpt fn run DIR/ --image gcr.io/example.com/my-fn

    # run the functions declared in files under FUNCTIONS_DIR/
    kpt fn run DIR/ --fn-path FUNCTIONS_DIR/

    # run the functions declared in files under DIR/
    kpt fn run DIR/
`

var RunShort = `Locally execute one or more functions`
var RunLong = `
Generate, transform, or validate configuration files using locally run containerized functions.

Programs are packaged as container images which are pulled and run locally.
If the container exits with non-zero status code, run will fail and print the
container ` + "`" + `STDERR` + "`" + `.

#### Imperatively run an single function

A function may be explicitly specified using the ` + "`" + `--image` + "`" + ` flag.

__Example:__ Locally run the container image ` + "`" + `gcr.io/example.com/my-fn` + "`" + ` against
the Resources in ` + "`" + `DIR/` + "`" + `:

	kpt fn run DIR/ --image gcr.io/example.com/my-fn

If ` + "`" + `DIR/` + "`" + ` is not specified, the source will default to STDIN and sink will default
to STDOUT.

__Example:__ this is equivalent to the preceding example:

	kpt source DIR/ |
	kpt fn run --image gcr.io/example.com/my-fn |
	kpt sink DIR/

Arguments specified after ` + "`" + `--` + "`" + ` will be provided to the function as a ` + "`" + `ConfigMap` + "`" + ` input.

__Example:__ In addition to the input Resources, provide to the container image a
` + "`" + `ConfigMap` + "`" + ` containing ` + "`" + `data: {foo: bar}` + "`" + `. This is used to parameterize the behavior
of the function:

	kpt fn run DIR/ --image gcr.io/example.com/my-fn -- foo=bar

#### Declaratively run one or more functions

Functions and their input configuration may be declared in files rather than directly
on the command line.

__Example:__ This is equivalent to the preceding example:

Create a file e.g. ` + "`" + `DIR/my-function.yaml` + "`" + `:

	apiVersion: v1
	kind: ConfigMap
	metdata:
	  annotations:
	    config.kubernetes.io/function: |
	      container:
	        image: gcr.io/example.com/my-fn
	data:
	  foo: bar

Run the function:

	kpt fn run DIR/

Here, rather than specifying ` + "`" + `gcr.io/example.com/my-fn` + "`" + ` as a flag, we specify it in a
file using the ` + "`" + `config.kubernetes.io/function` + "`" + ` annotation.

##### Scoping Rules

Functions which are nested under some sub directory are scoped only to Resources under 
that same sub directory. This allows fine grain control over how functions are 
executed:

__Example:__ Function declared in ` + "`" + `stuff/my-function.yaml` + "`" + ` is scoped to Resources in 
` + "`" + `stuff/` + "`" + ` and is NOT scoped to Resources in ` + "`" + `apps/` + "`" + `:

	.
	├── stuff
	│   ├── inscope-deployment.yaml
	│   ├── stuff2
	│   │     └── inscope-deployment.yaml
	│   └── my-function.yaml # functions is scoped to stuff/...
	└── apps
	    ├── not-inscope-deployment.yaml
	    └── not-inscope-service.yaml

Alternatively, you can also place function configurations in a special directory named 
` + "`" + `functions` + "`" + `.

__Example__: This is equivalent to previous example:

	.
	├── stuff
	│   ├── inscope-deployment.yaml
	│   ├── stuff2
	│   │     └── inscope-deployment.yaml
	│   └── functions
	│         └── my-function.yaml
	└── apps
	    ├── not-inscope-deployment.yaml
	    └── not-inscope-service.yaml

Alternatively, you can also use ` + "`" + `--fn-path` + "`" + ` to explicitly provide the directory 
containing function configurations:

	kpt fn run DIR/ --fn-path FUNCTIONS_DIR/

Alternatively, scoping can be disabled using ` + "`" + `--global-scope` + "`" + ` flag.

##### Declaring Multiple Functions

You may declare multiple functions. If they are specified in the same file 
(multi-object YAML file separated by ` + "`" + `---` + "`" + `), they will
be run sequentially in the order that they are specified.

##### Custom ` + "`" + `functionConfig` + "`" + `

Functions may define their own API input types - these may be client-side equivalents 
of CRDs:

__Example__: Declare two functions in ` + "`" + `DIR/functions/my-functions.yaml` + "`" + `:

	apiVersion: v1
	kind: ConfigMap
	metdata:
	  annotations:
	    config.kubernetes.io/function: |
	      container:
	        image: gcr.io/example.com/my-fn
	data:
	  foo: bar
	---
	apiVersion: v1
	kind: MyType
	metdata:
	  annotations:
	    config.kubernetes.io/function: |
	      container:
	        image: gcr.io/example.com/other-fn
	spec:
	  field:
	    nestedField: value
`
var RunExamples = `
	# read the Resources from DIR, provide them to a container my-fun as input,
	# write my-fn output back to DIR
	kpt fn run DIR/ --image gcr.io/example.com/my-fn
	
	# provide the my-fn with an input ConfigMap containing ` + "`" + `data: {foo: bar}` + "`" + `
	kpt fn run DIR/ --image gcr.io/example.com/my-fn:v1.0.0 -- foo=bar
	
	# run the functions in FUNCTIONS_DIR against the Resources in DIR
	kpt fn run DIR/ --fn-path FUNCTIONS_DIR/
	
	# discover functions in DIR and run them against Resource in DIR.
	# functions may be scoped to a subset of Resources -- see ` + "`" + `kpt help fn run` + "`" + `
	kpt fn run DIR/
`

var SinkShort = `Explicitly specify an output sink`
var SinkLong = `
Implements a Sink by reading command stdin and writing to a local directory.

    kpt fn sink [DIR]

  DIR:
    Path to local directory.  If unspecified, sink will write to stdout as if it were a single file.

` + "`" + `sink` + "`" + ` writes its input to a directory
`
var SinkExamples = `
    # run a function using explicit sources and sinks
    kpt fn source DIR/ | kpt run --image gcr.io/example.com/my-fn | kpt fn sink DIR/`

var SourceShort = `Explicitly specify an input source`
var SourceLong = `
Implements a Source by reading configuration and writing to command stdout.

    kpt fn source [DIR...]

  DIR:
    One or more paths to local directories.  Contents from directories will be concatenated.
    If no directories are provided, source will read from stdin as if it were a single file.

` + "`" + `source` + "`" + ` emits configuration to act as input to a function
`
var SourceExamples = `
    # print to stdout configuration formatted as an input source
    kpt fn source DIR/

    # run a function using explicit sources and sinks
    kpt fn source DIR/ | kpt run --image gcr.io/example.com/my-fn | kpt fn sink DIR/`
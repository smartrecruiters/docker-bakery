<p align="center">
	<h1 align="center">docker-bakery</h1>
	<p align="center">
		<a href="https://travis-ci.org/smartrecruiters/docker-bakery"><img alt="Build" src="https://travis-ci.org/smartrecruiters/docker-bakery.svg?branch=master"></a>	
		<a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>	
		<a href="https://goreportcard.com/report/github.com/smartrecruiters/docker-bakery"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/smartrecruiters/docker-bakery?style=flat-square"></a>
		<a href="http://godoc.org/github.com/smartrecruiters/docker-bakery"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-brightgreen.svg?style=flat-square"></a>
	</p>
</p>
<!-- MarkdownTOC  depth="4" autolink="true" bracket="round" autoanchor="true" -->

- [Purpose](#purpose)
- [Features](#features)
- [Example usage](#example-usage)
- [Example project](#example-project)
- [Assumptions](#assumptions)
- [Config](#config)
  - [Properties config section](#properties-config-section)
  - [Commands config section](#commands-config-section)
  - [Other config attributes](#other-config)
- [Dockerfile.template](#dockerfiletemplate)
- [Usage](#usage)
  - [Command help](#command-help)
  - [Command fill-template](#command-fill-template)
  - [Command build](#command-build)
  - [Command push](#command-push)
- [How to apply it to your project](#how-to-apply-it-to-your-project)
- [Limitations](#limitations)

<!-- /MarkdownTOC -->

<a id="purpose"></a>
# Purpose
Aim of the `docker-bakery` is to provide simple solution for automatic rebuilding of dependent images when parent image changes. 


<a id="features"></a>
# Features
- Automatic triggering of dependant images builds when parent changes
- Support for Dockerfile templating with usage of [golang template engine](https://golang.org/pkg/text/template/)
- Support for [semantic versioning](https://semver.org) scope changes
- Possibility to `build` and `push` docker images to custom registries
- Possibility of providing custom `build` and `push` commands
- Versioning with `git` tags
- Written in `golang`

<a id="example-usage"></a>
# Example usage
!["Example usage"](docker-bakery-demo.gif)

<a id="example-project"></a>
# Example project
See [docker-bakery-example](https://github.com/smartrecruiters/docker-bakery-example) to check how this tool works in action.

<a id="assumptions"></a>
# Assumptions
 - convention: image name is equal to the parent directory name
 - presence of Dockerfile.template will qualify image for automatic updates triggered by base images (`FROM` clause is analyzed)
 - first line in the dockerfile must start with `FROM ` otherwise dependency will not be determined for that file
 - scope change in the base image is propagated to the child images
 - to benefit from dependency updates the child image must use in template variable according to the convention: `{{.BASE_IMAGE_NAME_VERSION}}` where `BASE_IMAGE_NAME` should be substituted with uppercase directory name of the base image. 
 
 	Any hyphens (`-`) or dots (`.`) should be replaced with the underscores (`_`) in the variable name. 
 	
 	For example the correct variable name for image/directory named `jdk8-gradle2.14` is `{{.JDK8_GRADLE2_14_VERSION}}`  
 - docker file templates need to be placed in `git` repository (with defined remote) in order versioning of the images could work (versioning is done via `git tags`)
 - additional dynamic variables will be accessible for build templating
       
        
<a id="config"></a>
# Config
`docker-bakery` needs `config.json` to be provided. Config is used for templating `Dockerfile.template` files and build commands
Structure of the `config.json` file is as follows:

```
 {
	"properties": {
		"DEFAULT_PULL_REGISTRY": "some-private-registry.com:9084",
		"DEFAULT_PUSH_REGISTRY": "some-private-registry.com:9082",
		"DOCKERFILE_DIR": "Reserved dynamic property, contain path to currently build image. Can be used in template.",
		"IMAGE_NAME": "Reserved dynamic property, represents image name. Can be used in template.",
		"IMAGE_VERSION": "Reserved dynamic property, represents new version of the image. Can be used in template."
	},
	"commands": {
		"defaultBuildCommand": "docker build --tag {{.IMAGE_NAME}}:{{.IMAGE_VERSION}} --tag {{.DEFAULT_PUSH_REGISTRY}}/{{.IMAGE_NAME}}:{{.IMAGE_VERSION}} --tag {{.DEFAULT_PULL_REGISTRY}}/{{.IMAGE_NAME}}:{{.IMAGE_VERSION}} {{.DOCKERFILE_DIR}}",
		"defaultPushCommand": "docker push {{.DEFAULT_PUSH_REGISTRY}}/{{.IMAGE_NAME}}:{{.IMAGE_VERSION}}"
	},
	"reportFileName": "custom-report-filename.json",
	"verbose": false,
	"autoBuildExcludes": [
		"some-image-name-that-will-be-excluded-from-build-when-parent-changes"
	]
 }
```
 
<a id="properties-config-section"></a>
## Properties config section
 This section is dedicated for storing any custom properties that may be available for usage in `Dockerfile.template` files. 
 Feel free to modify this section and provide properties according to your needs. Flat structure should be preserved.
 
 This section will also be updated with dynamic properties during runtime. Dynamic properties do not have to be defined 
 in config as they are automatically added during runtime.
 
 Following properties belong to dynamic ones:
 - `DOCKERFILE_DIR` - will be replaced with currently processed dockerfile dir
 - `IMAGE_NAME` - will be replaced with currently processed image name
 - `IMAGE_VERSION` - will be replaced with currently processed image version
 - `*_VERSION` - where `*` is the image name. There will be that many properties of this kind as many images are in hierarchy. Initially those properties will be filled with latest versions of pushed images.
 - `BAKERY_BUILDER_NAME` - will be replaced with the git user name (taken from `git config user.name`)  
 - `BAKERY_BUILDER_EMAIL` - will be replaced with the git user email (taken from `git config user.email`)
 - `BAKERY_BUILDER_HOST` - will be replaced with hostname of the machine where build is executed
 - `BAKERY_BUILD_DATE` - will be replaced with current build date 
 - `BAKERY_IMAGE_HIERARCHY` - will be replaced with the path representing image hierarchy in the following format: 
 
 `parent1:versionOfParent1->parent2:versionOfParent2->imageName:imageVersion` 
 
 Hierarchy is built automatically given that parent images are exporting the same `ENV` variable that can be accessed in child images. Check the example project for references.
 - `BAKERY_SIGNATURE_VALUE` - will be replaced with a one liner string value embedding other `BAKERY*` variables together. Can be used in templates to create for example `ENV` variable. Example:
  
  `SINGATURE=Builder Name;builder@email.com;builder-host-name;2018-03-16 15:47:58;alpine-java:8u144b01_jdk->mammal:3.2.0->dog:4.0.0->dobermann:4.0.0->smaller-dobermann:4.0.0` 
 - `BAKERY_SIGNATURE_ENVS` - will be replaced with embedded `BAKERY*` variables in a `key=value` format. Convenient if you wish to have all `BAKERY*` variables in a dockerfile under single key. 
 Check the example project for references. 

<a id="commands-config-section"></a>
## Commands config section
This section contains two templates used for building and pushing docker images. It allows for specifying custom parameters. 
Commands defined here as templates will be filled with available defined properties from the config section + the dynamic properties set during runtime. 

<a id="other-config"></a>
## Other config attributes
  `reportFileName` - if set it will be used as a file name to store information (in JSON format) about successfully built images. 

<a id="dockerfiletemplate"></a>
# Dockerfile.template
Presence of the `Dockerfile.template` file qualifies the image for the place in hierarchy and therefore allows for triggering builds that depend from this image. It also ensures that image build will be triggered when its parent changes. 

<a id="usage"></a>
# Usage
To make use of `docker-bakery` as convenient as possible checkout usage of `Makefiles` from the [example project](https://github.com/smartrecruiters/docker-bakery-example) that will simplify usage greatly.
If you don't want to use makefiles you can still use `docker-bakery` tool directly.
Checkout the CLI help via `docker-bakery -h`. 

<a id="command-help"></a>
## Command help
```
COMMANDS:
     fill-template, prepare, prepare-recipe         Used to fill Dockerfile.template file. Values needed for template are taken from the config file and from dynamic properties provided during runtime.
     build                                          Used to build next version of the images in given scope. Optionally it can skip build of dependant images.
     push                                           Used to push next version of the images in given scope. Optionally it can skip push of dependant images.
     show-structure, ss, show-hierarchy, hierarchy  Used to display hierarchy of the images
     dump-latest-versions, dump                     Used to dump data about latest versions of images to the provided file
     help, h                                        Shows a list of commands or help for one command

```
<a id="command-fill-template"></a>
## Command fill-template
```
docker-bakery fill-template -h
NAME:
   docker-bakery fill-template - Used to fill Dockerfile.template file. Values needed for template are taken from the config file and from dynamic properties provided during runtime.

USAGE:
   docker-bakery fill-template [command options] [arguments...]

OPTIONS:
   --input value, -i value      Required. Input Dockerfile.template.
   --output value, -o value     Required. Output Dockerfile generated from template.
   --config value, -c value     Required. Path to config.json with properties and build commands defined.
   --rootDir value, --rd value  Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.
```
<a id="command-build"></a>
## Command build
```
docker-bakery build -h
NAME:
   docker-bakery build - Used to build next version of the images in given scope. Optionally it can skip build of dependant images.

USAGE:
   docker-bakery build [command options] [arguments...]

OPTIONS:
   --dockerfile value, -d value  Required. Path to dockerfile/dockerfile.template file that needs to be build.
   --scope value, -s value       Required. Scope of the change used to generate the next version. Can be one of: major/minor/patch.
   --config value, -c value      Required. Path to config.json with properties and build commands defined.
   --root-dir value, --rd value  Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.
   --skip-dependants, --sd       Optional. False be default. If this flag is set build of the parent will not trigger dependant builds.
   --property value, -p value    Optional. Allows for providing additional multiple properties that can be used during templating. Overrides properties defined in config.json file. Expected format is: -p propertyName=propertyValue
     
```
<a id="command-push"></a>
## Command push
```
docker-bakery push -h
NAME:
   docker-bakery push - Used to push next version of the images in given scope. Optionally it can skip push of dependant images.

USAGE:
   docker-bakery push [command options] [arguments...]

OPTIONS:
   --dockerfile value, -d value  Required. Path to the dockerfile/dockerfile.template that needs to be pushed.
   --scope value, -s value       Required. Scope of the change used to generate the next version. Can be one of: major/minor/patch.
   --config value, -c value      Required. Path to config.json with properties and build commands defined.
   --rootDir value, --rd value   Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.
   --skip-dependants, --sd       Optional. False be default. If this flag is set build of the parent will not trigger dependant builds.

```

<a id="how-to-apply-it-to-your-project"></a>
# How to apply it to your project
Applying `docker-bakery` is quite simple. Take a look [here](https://github.com/smartrecruiters/docker-bakery-example#how-to-apply-it-to-your-project)

<a id="limitations"></a>
# Limitations
At the moment docker multi-stage builds are not fully supported. Only the first line `FROM` is taken into consideration when determining the parent of the image and its place in the hierarchy. 

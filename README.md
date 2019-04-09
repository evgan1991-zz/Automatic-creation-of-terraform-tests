## Description

This script is designed to automatically generate a simple test pattern for [Terraform](https://www.terraform.io/) modules using [Terratest](https://github.com/gruntwork-io/terratest). The script automatically creates additional folders and adds the files necessary for the tests. The only test added by this script checks the possibility to create and destroy a resource created by your terraform module.


## Usage

For creating test
```hcl
go run creator.go ~./path/to/your/terraform/module
```

For running test
```hcl
go test ~./path/to/your/terraform/module/test
```


## Files to be added
 * examples/main.tf
 * examples/variables.tf
 * examples/outputs.tf
 * [test/automatically_generated_test.go](templates/test)


## Files to be changed
 * [.travis.yml](templates/travis)
 * [.gitignore](templates/gitignore)

From the links you can see what will be added to these files. If you apply the script several times, this part will be added to .travis and .gitignore several times.


## [Automatically_generated_test.go](templates/test)
In order to avoid a name conflict during execution, a random line of text is added to the name of the resource. By default, the name is set by the parameter "name" in the description of the structure that is passed to the module as input parameters.

```hcl
Vars: map[string]interface{}{
    "aws_region": region,
    "name"      : "test_name_" + randSeq(10),
},
```
If the variable name or any other identifier of your resource has a different name, change this name in [automatically_generated_test](templates/test)


## Terraform versions

Terraform version 0.11.11 or newer is required for this module to work.


## Go versions

go version go1.12 darwin/amd64


## Authors

Yauheni Anashkin

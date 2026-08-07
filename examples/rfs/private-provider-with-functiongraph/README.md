# Create a Private Provider with FunctionGraph Backend

This example provides best practice code for using Terraform to create an RFS private provider in HuaweiCloud Resource
Formation Service (RFS). The example first creates a FunctionGraph function as the execution backend, then creates a
private provider associated with the function, and finally publishes a new version for the private provider.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RFS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `function_name` - The name of the FunctionGraph function
* `function_code` - The inline code content of the FunctionGraph function
* `private_provider_name` - The name of the RFS private provider

#### Optional Variables

* `function_app` - The group name of the FunctionGraph function (default: "default")
* `function_handler` - The handler of the FunctionGraph function (default: "index.handler")
* `function_memory_size` - The memory size of the FunctionGraph function in MB (default: 128)
* `function_timeout` - The timeout of the FunctionGraph function in seconds (default: 3)
* `function_runtime` - The runtime of the FunctionGraph function (default: "Node.js12.13")
* `private_provider_description` - The description of the RFS private provider (default: `""`)
* `private_provider_version` - The initial version number of the RFS private provider (default: "1.0.0")
* `private_provider_version_description` - The initial version description of the RFS private provider (default: `""`)
* `provider_version_number` - The version number of the RFS private provider version (default: "2.0.0")
* `provider_version_description` - The description of the RFS private provider version (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  function_name         = "tf-test-function"
  function_code         = <<-EOT
  exports.handler = async (event, context) => {
      const result =
      {
          'statusCode': 200,
          'headers':
          {
              'Content-Type': 'application/json'
          },
          'isBase64Encoded': false,
          'body': JSON.stringify(event)
      }
      return result
  }
  EOT
  private_provider_name = "tf-test-provider"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The `private_provider_name` only allows lowercase letters, digits, and hyphens (-), and must be unique within its
  domain and region
* The FunctionGraph function must be in the same region as the RFS resources
* The `function_graph_urn`, `provider_version` and `version_description` of the private provider are non-updatable,
  changing them will recreate the private provider
* All arguments of the private provider version are non-updatable, changing them will recreate the version
* Deleting the private provider will also delete all its versions

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.91.0 |

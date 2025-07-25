# Register an API with Custom Authorizer and FunctionGraph

This example provides best practice code for using Terraform to create an API Gateway instance with a custom authorizer
and FunctionGraph backend on HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `function_name` - The name of the FunctionGraph function
* `function_code` - The code content of the FunctionGraph function
* `instance_name` - The instance name of the dedicated APIG
* `enterprise_project_id` - The ID of the enterprise project to which the APIG instance belongs (required for
  enterprise users, default: null)
* `custom_authorizer_name` - The name of the custom authorizer
* `response_name` - The response name of the dedicated APIG
* `response_rules` - The response rules of the dedicated APIG
  - `error_type` - The error type of the API response rule
  - `body` - The body template of the API response rule
  - `status_code` - The HTTP status code of the API response rule
  - `headers` - The configuration of the custom response headers (default: [])
    + `key` - The key name of the response header
    + `value` - The value for the specified response header key
* `group_name` - The group name of the dedicated APIG
* `api_name` - The name of the API
* `api_request_path` - The request path of the API
* `api_backend_params` - The backend parameters of the API
  - `type` - The type of the backend parameter
  - `name` - The name of the backend parameter
  - `location` - The location of the backend parameter
  - `value` - The value of the backend parameter
  - `system_param_type`- The system parameter type of the backend parameter (required for system parameters)

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `function_memory_size` - The memory size (MB) for the function (default: 128)
* `function_runtime` - The runtime environment for the function (default: "Python3.9")
* `function_timeout` - The timeout (seconds) for the function (default: 3)
* `function_handler` - The handler of the function (default: "index.handler")
* `function_code_type` - The code type of the function (default: "inline")
* `function_app` - The name of the application to which the function belongs (default: "default")
* `availability_zones` - The availability zones to which the APIG instance belongs (default: [])
  If not specified, will be automatically allocated based on the number of availability_zones_count
* `instance_edition` - The edition of the APIG instance (default: "BASIC")
* `availability_zones_count` - The number of availability zones to which the instance belongs (default: 1)
* `function_version` - The version of the function (default: "latest")
* `custom_authorizer_type` - The type of the custom authorizer (default: "FRONTEND")
* `custom_authorizer_network_type` - The network type of the custom authorizer (default: "V1")
* `custom_authorizer_cache_age` - The cache age of the custom authorizer (default: 0)
* `custom_authorizer_is_body_send` - Whether to send body in the custom authorizer (default: false)
* `custom_authorizer_use_data` - The user data used to obtain backend access authorization (default: null)
* `custom_authorizer_identity` - The identity list of the custom authorizer (default: [])
  - `name` - The name of the identity
  - `location` - The location of the identity
  - `validation` - The parameter verification expression
* `api_type` - The type of the API (default: "Public")
* `api_request_protocol` - The request protocol of the API (default: "BOTH")
* `api_request_method` - The request method of the API (default: "GET")
* `api_matching` - The matching rule of the API (default: "Exact")
* `api_func_graph_network_type` - The network type of the FunctionGraph function (default: "V1")
* `api_func_graph_request_protocol` - The request protocol of the FunctionGraph function (default: "HTTPS")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables. Example:

  ```hcl
  vpc_name               = "your_vpc_name"
  subnet_name            = "your_subnet_name"
  security_group_name    = "your_security_group_name"
  function_name          = "your_function_name"
  function_code          = "your_function_code"
  instance_name          = "your_apig_instance_name"
  enterprise_project_id  = "your_enterprise_project_id"
  custom_authorizer_name = "your_custom_authorizer_name"
  group_name             = "your_group_name"
  response_name          = "your_response_name"
  response_rules         = "your_response_rules"
  api_name               = "your_api_name"
  api_request_path       = "your_api_request_path"
  api_backend_params     = "your_api_backend_params"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.
* Please be aware of API Gateway and FunctionGraph service quotas in your HuaweiCloud account.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.72.0 |

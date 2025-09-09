# Create FunctionGraph CTS Trigger

This example provides best practice code for using Terraform to create a FunctionGraph CTS (Cloud Trace Service) trigger
in HuaweiCloud FunctionGraph service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* FunctionGraph service enabled in the target region
* CTS service enabled in the target region
* Authorize FunctionGraph to use the CTS service as an administrator

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the FunctionGraph service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `function_name` - The name of the FunctionGraph function
* `function_agency_name` - The agency name of the FunctionGraph function
* `trigger_name` - The name of the FunctionGraph CTS trigger
* `trigger_operations` - The operations to monitor for the FunctionGraph CTS trigger

#### Optional Variables

* `function_memory_size` - The memory size of the function in MB (default: 128)
* `function_timeout` - The timeout of the function in seconds (default: 10)
* `function_runtime` - The runtime of the function (default: "Python2.7")
* `function_code` - The source code of the function (default: provided Python example)
* `function_description` - The description of the function (default: "")
* `trigger_status` - The status of the FunctionGraph CTS trigger (default: "ACTIVE")

### CTS Operations

Each CTS operation format is as follows:
  `service:resource:operation1;operation2;operation3;...`.

For example:

* `SMN:subscription:addSubscriptionFromSubscriptionUser;subscribe`
* `FunctionGraph:Functions:createFunction`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  function_name        = "tf_test_cts_function"
  function_agency_name = "function_all_trust"
  trigger_name         = "tf_test_cts_trigger"
  trigger_operations   = [
    "FunctionGraph:Functions:createFunction"
  ]
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

## Example Function Code

The default function code is a simple Python function that logs the CTS trigger event:

```python
# -*- coding:utf-8 -*-
'''
CTS trigger event:
{
  "cts":  {
        "time": "",
        "user": {
            "name": "userName",
            "id": "",
            "domain": {
                "name": "domainName",
                "id": ""
            }
        },
        "request": {},
        "response": {},
        "code": 204,
        "service_type": "FunctionGraph",
        "resource_type": "",
        "resource_name": "",
        "resource_id": {},
        "trace_name": "",
        "trace_type": "ConsoleAction",
        "record_time": "",
        "trace_id": "",
        "trace_status": "normal"
    }
}
'''
def handler (event, context):
    trace_name = event["cts"]["resource_name"]
    timeinfo = event["cts"]["time"]
    print(timeinfo+' '+trace_name)
```

You can customize the function code by modifying the `function_code` variable in your `terraform.tfvars` file.

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The CTS trigger is dependent on the FunctionGraph function
* Please read the implicit and explicit dependencies in the script carefully

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.11.0 |
| huaweicloud | >= 1.62.1 |

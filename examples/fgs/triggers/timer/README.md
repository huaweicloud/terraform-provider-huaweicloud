# Create FunctionGraph Timer Trigger (Cron)

This example provides best practice code for using Terraform to create a FunctionGraph timer trigger with cron schedule
type in HuaweiCloud FunctionGraph service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* FunctionGraph service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the FunctionGraph service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `function_name` - The name of the FunctionGraph function
* `trigger_name` - The name of the FunctionGraph timer trigger

#### Optional Variables

* `function_memory_size` - The memory size of the function in MB (default: 128)
* `function_timeout` - The timeout of the function in seconds (default: 10)
* `function_runtime` - The runtime of the function (default: "Python2.7")
* `function_code` - The source code of the function (default: provided Python example)
* `function_description` - The description of the function (default: "")
* `trigger_status` - The status of the FunctionGraph timer trigger (default: "ACTIVE")
* `trigger_schedule_type` - The schedule type of the FunctionGraph timer trigger (default: "Cron")
* `trigger_sync_execution` - Whether to execute the function synchronously (default: false)
* `trigger_user_event` - The user event description for the FunctionGraph timer trigger (default: "")
* `trigger_schedule` - The schedule expression for the FunctionGraph timer trigger (default: "@every 1h30m")

### Cron expressions

Cron expressions can be learned by referring to the FunctionGraph [documentation](https://support.huaweicloud.com/usermanual-functiongraph/functiongraph_01_0207.html#functiongraph_01_0207__section1237972211317).

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  function_name        = "tf_test_timer_function"
  function_description = "Created by Terraform for timer trigger example"
  trigger_name         = "tf_test_timer_cron"
  trigger_user_event   = "Timer trigger with Cron schedule type, triggered every three days"
  trigger_schedule     = "@every 3d"
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

The default function code is a simple Python function that logs the trigger event:

```python
import json

def handler(event, context):
    print("Timer trigger executed!")
    print("Event:", json.dumps(event))
    print("Context:", json.dumps(context))
    return {
        'statusCode': 200,
        'body': json.dumps('Hello, FunctionGraph!')
    }
```

You can customize the function code by modifying the `function_code` variable in your `terraform.tfvars` file.

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The timer trigger is dependent on the FunctionGraph function
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* Timer trigger supports both ACTIVE and DISABLED status
* The function supports various runtimes including Python, Node.js, Java, Go, PHP, and C#
* Cron schedule supports both standard cron expressions and simplified expressions like `@every 1h30m`

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.62.1 |

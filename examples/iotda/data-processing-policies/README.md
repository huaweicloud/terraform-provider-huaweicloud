# Create IoTDA Data Flow Control and Backlog Policies

This example provides best practice code for using Terraform to create an IoTDA data flow
control policy and a data backlog policy for data forwarding in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An IoTDA **standard** or **enterprise** edition instance has been created in the console
  (these two kinds of policies are only supported on standard and enterprise edition
  instances, the free unit `S0` is sufficient and does not incur extra fees)
* The HTTPS application access address of the instance has been obtained: login to the IoTDA
  console, choose the instance **Overview** and click **Access Details**

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IoTDA service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `flow_control_policy_name` - The name of the data flow control policy
* `flow_control_policy_limit` - The flow control limit of the policy in tps
* `backlog_policy_name` - The name of the data backlog policy
* `backlog_policy_size` - The size of data backlog in bytes
* `backlog_policy_time` - The data backlog time in seconds
* `iotda_access_address` - The HTTPS application access address of the IoTDA instance

#### Optional Variables

* `flow_control_policy_description` - The description of the data flow control policy
* `backlog_policy_description` - The description of the data backlog policy

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  flow_control_policy_name        = "tf_test_iotda_flow_control_policy"
  flow_control_policy_description = "Limit-the-data-forwarding-tps-of-the-tenant"
  flow_control_policy_limit       = 500
  backlog_policy_name             = "tf_test_iotda_backlog_policy"
  backlog_policy_description      = "Control-the-size-and-time-of-forwarded-data-backlog"
  backlog_policy_size             = "524288000"
  backlog_policy_time             = "3600"
  iotda_access_address            = "https://9bc34xxxxx.st1.iotda-app.cn-north-4.myhuaweicloud.com"
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

* These two kinds of policies are only supported on IoTDA standard and enterprise edition
  instances, the instance access address must be specified through the `endpoints.iotda`
  argument of the provider block, and the provider version must be `1.63.0` or later which
  supports the derived authentication required by these instances
* The `scope` argument of the flow control policy is fixed to **USER** (tenant level) in
  this example so the `scope_value` argument is not required, and both arguments are
  `ForceNew`, changing them will recreate the policy
* The `limit` argument must be an integer between `1` and `1,000` (tps), the default
  value is `1,000`
* The valid range of the `backlog_size` argument is integers from `0` to `1,073,741,823`
  bytes, and the valid range of the `backlog_time` argument is integers from `0` to
  `86,399` seconds, `0` means no backlog, when both dimensions are configured the
  dimension that reaches the threshold first shall prevail
* The name and description of the policies must not exceed `256` characters, only Chinese
  characters, letters, digits and the following characters are allowed: `_?'#().,&%@!-`,
  spaces are not allowed
* If the `apply` command reports `Resource not found` or `inconsistent result after
  apply` immediately after a resource is created, it is caused by the data synchronization
  delay of the service, wait a few seconds and run `terraform apply` again
* Both resources support import using the `id`

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.71.0 |

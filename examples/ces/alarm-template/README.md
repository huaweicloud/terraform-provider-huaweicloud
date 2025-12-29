# Create a CES alarm template example

This example provides best practice code for using Terraform to create an alarm template in HuaweiCloud CES service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The name of the CES alarm template.
* `policies` - The policy list of the CES alarm template.
  + `namespace` - The namespace of the service.
  + `metric_name` - The alarm metric name.
  + `period` - The judgment period of alarm condition.
  + `filter` - The data rollup methods.
  + `comparison_operator` - The comparison conditions for alarm threshold.
  + `count` - The number of consecutive alarm triggering times.
  + `suppress_duration` - The alarm suppression cycle.
  + `value` - The alarm threshold.
  + `hierarchical_value` - The multiple levels of alarm thresholds.
    - `critical` - The threshold for the critical level.
    - `major` - The threshold for the major level.
    - `minor` - The threshold for the minor level.
    - `info` - The threshold for the info level.
  + `alarm_level` - The alarm level.
  + `unit` - The unit string of the alarm threshold.
  + `dimension_name` - The resource dimension.

#### Optional Variables

* `type` - The type of the CES alarm template.
* `description` - The description of the CES alarm template.
* `is_overwrite` - Whether to overwrite an existing alarm template with the same template name.
* `delete_associate_alarm` -  Whether delete the alarm rule which the alarm template associated with.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  name     = "tf_test_ces_alarm_template_name"
  policies = [
    {
      namespace           = "SYS.APIG"
      dimension_name      = "api_id"
      metric_name         = "req_count_2xx"
      period              = 1
      filter              = "average"
      comparison_operator = "="
      value               = "10"
      unit                = "times/minute"
      count               = 3
      alarm_level         = 2
      suppress_duration   = 300
      hierarchical_value  = [
        {
           critical = 12
           major    = 10
           minor    = 8
           info     = 4
        }
      ]
    }
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

## Note

* Make sure to keep your credentials secure and never commit them to version control

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.77.1 |

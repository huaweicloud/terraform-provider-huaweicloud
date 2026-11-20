# Create an AS bandwidth Policy

This example provides best practice code for using Terraform to create an Auto Scaling (AS) bandwidth policy in
HuaweiCloud. The example demonstrates how to set up a complete AS environment with bandwidth scaling policy.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* AS service enabled in the target region
* CES service enabled in the target region
* EIP service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bandwidth_name` - The name of the bandwidth name
* `alarm_rule_name` - The name of the CES alarm rule
* `rule_conditions` - The conditions of the alarm rule
  + `period` - The period for collecting the metric data in seconds
  + `filter` - The data aggregation method (average, max, min, sum)
  + `comparison_operator` - The comparison operator (>, <, >=, <=, =)
  + `value` - The threshold value for the alarm
  + `count` - The number of consecutive periods that the condition must be met
  + `unit` - The data unit For details. see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).
  + `alarm_level` - The alarm level (1-4, default: 2)
  + `suppress_duration` - The suppression duration in seconds (default: 0)
* `bandwidth_policy_name` - The name of the AS bandwidth policy
* `interval_alarm_conditions` - The name of the scaling down policy
  + `lower_bound` - The lower limit of the value range. (default: null)
  + `upper_bound` - The upper limit of the value range. (default: null)
  + `operation` - The operation to be performed (ADD,REDUCE,SET)
  + `size` - The the operation size (1-300, default: 2)
  + `limits` - The operation restrictions (1-2,000)
    If `operation` is not **SET**, this parameter takes effect.

#### Optional Variables

* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_type` - The type of the bandwidth (default: "share")
* `enterprise_project_id` - The enterprise project ID (default: null)
* `charging_mode` - The charging mode of the bandwidth (default: "postPaid")
* `period_unit` - The period unit of the bandwidth, only required when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the bandwidth, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - Whether the auto renew is enabled (default: "false")
* `alarm_rule_description` - The description of the CES alarm rule (default: "")
* `alarm_enabled` - Whether the CES alarm rule is enabled (default: true)
* `alarm_action_enabled` - Whether to enable the action to be triggered by an alarm (default: true)
* `bandwidth_policy_type` - The type of the AS bandwidth policy (default: "INTERVAL_ALARM")
* `bandwidth_policy_description` - The description of the AS bandwidth policy (default: "")
* `cool_down_time` - The cooldown period of the AS bandwidth policy (default: 300)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bandwidth_name  = "your_bandwidth_name"
  alarm_rule_name = "your_alarm_rule_name"

  rule_conditions = [
    {
      period              = 300
      filter              = "average"
      comparison_operator = ">"
      value               = 3
      count               = 5
      unit                = "Mbit/s"
    }
  ]

  bandwidth_policy_name = "your_bandwidth_policy_name"

  interval_alarm_conditions = [
    {
        lower_bound = "3"
        upper_bound = "10"
        operation   = "ADD"
        size        = 1
        limits      = 20
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

## Architecture

This example creates the following resources:

1. **EIP Bandwidth** - EIP bandwidth for AS bandwidth policy
2. **CES Alarm Rule** - Cloud monitoring alarm rule based on CPU utilization
3. **AS Bandwidth Policy** - Scaling bandwidth policy triggered by alarms

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.73.9 |

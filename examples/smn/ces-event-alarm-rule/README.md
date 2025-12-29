# Configure CES Event Alarm Rule for SMN

This example provides best practice code for using Terraform to configure a CES (Cloud Eye Service) alarm rule for
monitoring SMN (Simple Message Notification) metrics on HuaweiCloud. The example demonstrates how to create an SMN
topic and set up a CES alarm rule to monitor SMN metrics and send notifications when specified conditions are met.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `smn_topic_name` - The name of the SMN topic used to send alarm notifications
* `alarm_rule_name` - The name of the CES alarm rule
* `alarm_rule_conditions` - The list of alarm rule conditions.
  - `metric_name` (required) - The name of the metric to monitor
  - `period` (required) - The period of the alarm condition in seconds
  - `filter` (required) - The filter method (average, max, min, sum)
  - `comparison_operator` (required) - The comparison operator (>, >=, <, <=, =)
  - `value` (required) - The threshold value
  - `count` (required) - The number of consecutive periods that the alarm condition must be met
  - `unit` (optional) - The unit of the alarm condition value
  - `suppress_duration` (optional) - The suppression duration in seconds after an alarm is triggered
  - `alarm_level` (optional) - The alarm level (1: critical, 2: major, 3: minor, 4: informational)

#### Optional Variables

* `enterprise_project_id` - The ID of the enterprise project (default: null)
* `alarm_rule_description` - The description of the CES alarm rule (default: null)
* `alarm_action_enabled` - Whether to enable the action to be triggered by an alarm (default: true)
* `alarm_enabled` - Whether to enable the alarm (default: true)
* `alarm_type` - The type of the alarm. Valid values: EVENT.SYS, EVENT.CUSTOM, MULTI_INSTANCE, ALL_INSTANCE (default: "ALL_INSTANCE")
* `alarm_rule_resource` - The list of resource dimensions for specified monitoring scope (default: [])
  - `name` (required) - The dimension name
  - `value` (optional) - The dimension value
* `alarm_rule_notification_begin_time` - The alarm notification start time (default: null)
* `alarm_rule_notification_end_time` - The alarm notification stop time (default: null)
* `alarm_rule_effective_timezone` - The time zone (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  smn_topic_name         = "tf_test_topic"
  alarm_rule_name        = "tf_test_alarm_rule"
  alarm_rule_description = "Monitor SMN topic events"
  alarm_type             = "ALL_INSTANCE"
  alarm_rule_conditions  = [
    {
      metric_name         = "email_total_count"
      period              = "1"
      filter              = "average"
      comparison_operator = ">="
      value               = "80"
      count               = "3"
      unit                = "count"
      alarm_level         = "3"
    }
  ]

  alarm_rule_resource = [
    {
      name = "topic_id"
    }
  ]
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The metric namespace is fixed to "SYS.SMN" for SMN monitoring
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.77.1 |

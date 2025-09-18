# Create a LTS SQL alarm rule

This example provides best practice code for using Terraform to create a LTS SQL alarm rule in HuaweiCloud
LTS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* LTS service enabled in the target region
* SMN service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the LTS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `group_name` - The name of the log group
* `stream_name` - The name of the log stream
* `topic_name` - The name of the SMN topic
* `alarm_rule_name` - The name of the SQL alarm rule
* `alarm_rule_condition_expression` - The condition expression of the SQL alarm rule
* `alarm_rule_request_title` - The request title of the SQL alarm rule
* `alarm_rule_request_sql` - The request SQL of the SQL alarm rule
* `alarm_rule_notification_user_name` - The notification user name of the SQL alarm rule

#### Optional Variables

* `group_log_expiration_days` - The log expiration days of the log group (default: 14)
* `stream_log_expiration_days` - The log expiration days of the log stream (default: null)
* `notification_template_name` - The name of the notification template (default: "")
* `domain_id` - The domain ID (default: null)
* `alarm_rule_alarm_level` - The alarm level of the SQL alarm rule (default: "MINOR")
* `alarm_rule_trigger_condition_count` - The trigger condition count of the SQL alarm rule (default: 2)
* `alarm_rule_trigger_condition_frequency` - The trigger condition frequency of the SQL alarm rule (default: 3)
* `alarm_rule_send_recovery_notifications` - The send recovery notifications of the SQL alarm rule (default: true)
* `alarm_rule_recovery_frequency` - The recovery frequency of the SQL alarm rule (default: 4)
* `alarm_rule_notification_frequency` - The notification frequency of the SQL alarm rule (default: 15)
* `alarm_rule_alias` - The alias of the SQL alarm rule (default: "")
* `alarm_rule_request_search_time_range_unit` - The request search time range unit of the SQL alarm rule (default: "minute")
* `alarm_rule_request_search_time_range` - The request search time range of the SQL alarm rule (default: 5)
* `alarm_rule_frequency_type` - The frequency type of the SQL alarm rule (default: "HOURLY")
* `alarm_rule_notification_language` - The notification language of the SQL alarm rule (default: "en-us")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  group_name                        = "tf_test_log_group"
  stream_name                       = "tf_test_log_stream"
  topic_name                        = "tf-test-topic"
  domain_id                         = "your_domain_id"
  alarm_rule_name                   = "tf-test-sql-alarm-rule"
  alarm_rule_condition_expression   = "t>0"
  alarm_rule_request_title          = "tf-test-sql-alarm-rule-title"
  alarm_rule_request_sql            = "select count(*) as t"
  alarm_rule_notification_user_name = "your_notification_user_name"
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
* The SQL alarm rule is dependent on the log group, log stream, and SMN topic
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* The SQL alarm rule supports different alarm levels: CRITICAL, MAJOR, MINOR, INFO
* The frequency type supports: HOURLY, DAILY, WEEKLY, CRON, FIXED_RATE
* The SQL query statement needs to comply with LTS SQL syntax specifications
* The alarm condition expression variable names need to correspond to the field aliases in the SQL query

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.74.1 |

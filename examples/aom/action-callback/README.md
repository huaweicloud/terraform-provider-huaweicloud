# AOM Alarm Action Callback

This example provides best practice code for using Terraform to configure AOM alarm action callback through SMN (Simple
Message Notification) on HuaweiCloud. The example demonstrates how to create a custom alarm notification template and
configure alarm action rules with callback URLs to receive alarm notifications.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* AOM service enabled in the target region
* SMN service enabled in the target region
* A callback endpoint URL that can receive HTTP/HTTPS requests

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `smn_topic_name` - The name of the SMN topic used to send notifications
* `alarm_callback_urls` - The list of callback URLs to receive alarm notifications (must start with `http://` or
  `https://`, at least one URL is required)
* `alarm_notification_template_name` - The name of the AOM alarm notification template
* `alarm_action_rule_name` - The name of the AOM alarm action rule
* `alarm_action_rule_user_name` - The user name of the AOM alarm action rule

#### Optional Variables

* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `alarm_notification_template_locale` - The locale of the alarm notification template (default: "en-us", valid values:
  "en-us" or "zh-cn")
* `alarm_notification_template_description` - The description of the alarm notification template (default: "")
* `alarm_notification_template_notification_type` - The notification type of the template (default: "email")
* `alarm_notification_template_notification_topic` - The notification topic/title template (default: "An alert occurred
  at time $${starts_at}[$${event_severity}_$${event_type}_$${clear_type}].")
* `alarm_notification_template_content` - The content template of the alarm notification (default: includes alarm name,
  ID, trigger time, severity, content, resource identifier, and remediation suggestion)
* `alarm_action_rule_type` - The type of the AOM alarm action rule (default: "1" for notification)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  smn_topic_name                          = "tf_test_alarm_action_callback"
  alarm_callback_urls                     = ["https://www.example.com/alarm-callback"]
  alarm_notification_template_name        = "tf_test_alarm_action_callback"
  alarm_notification_template_description = "This is a AOM alarm notification template created by Terraform"
  alarm_action_rule_name                  = "tf_test_alarm_action_callback"
  alarm_action_rule_user_name             = "your_operation_user_name"
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
* The alarm callback URLs must be valid HTTP or HTTPS URLs and must be accessible from the internet
* The alarm action rule depends on the alarm notification template
* The SMN subscription will be created automatically for each callback URL in the list
* The protocol (HTTP/HTTPS) for SMN subscription is automatically detected from the URL prefix
* The notification template supports variable substitution using `$${variable_name}` syntax
* Available template variables include:
  + `$${event_name_alias}` - Alarm name alias
  + `$${id}` - Alarm ID
  + `$${action_rule}` - Notification rule name
  + `$${starts_at}` - Trigger time
  + `$${event_severity}` - Alarm severity level
  + `$${event_type}` - Event type
  + `$${clear_type}` - Clear type
  + `$${alarm_info}` - Alarm content
  + `$${resources_new}` - Resource identifier
  + `$${alarm_fix_suggestion_zh}` - Remediation suggestion (Chinese)
* The template locale determines the language of built-in template variables
* All resources will be created in the specified region
* The callback endpoint should be able to handle POST requests with JSON payload
* When an alarm is triggered, SMN will send a notification to the configured callback URLs

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | = 1.80.4 |

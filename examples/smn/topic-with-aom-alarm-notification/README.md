# Configure AOM Alarm Notification via SMN and LTS

This example provides best practice code for using Terraform to configure AOM (Application Operations Management) alarm
notifications through SMN (Simple Message Notification) and LTS (Log Tank Service) on HuaweiCloud. The example
demonstrates how to create an SMN topic, configure LTS group and stream, associate SMN topic with LTS for log storage,
and configure an AOM alarm action rule to send notifications through SMN.

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

* `smn_topic_name` - The name of the SMN topic used to send AOM alarm notifications
* `lts_group_name` - The name of the LTS group for storing alarm logs
* `lts_stream_name` - The name of the LTS stream for storing alarm logs
* `alarm_action_rule_name` - The name of the AOM alarm action rule
* `alarm_action_rule_user_name` - The user name of the AOM alarm action rule

#### Optional Variables

* `enterprise_project_id` - The ID of the enterprise project (default: null)
* `lts_group_ttl_in_days` - The TTL in days of the LTS group (default: 30)
* `alarm_action_rule_type` - The type of the AOM alarm action rule (default: "1" for notification)
* `alarm_action_rule_description` - The description of the AOM alarm action rule (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  smn_topic_name              = "tf_test_aom_notification"
  lts_group_name              = "tf_test_aom_notification"
  lts_stream_name             = "tf_test_aom_notification"
  alarm_action_rule_name      = "tf_test_aom_notification"
  alarm_action_rule_user_name = "your_operation_user_name"
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
* The SMN topic will be associated with LTS for storing alarm notification logs
* The AOM alarm action rule uses the built-in notification template "aom.built-in.template.zh" (Chinese)
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.4 |

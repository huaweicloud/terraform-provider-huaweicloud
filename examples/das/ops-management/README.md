# Configure DAS operations management including instance group and email notification

This example provides best practice code for using Terraform to configure operations management in HuaweiCloud
DAS service, including instance group management, email template configuration, and batch email sending.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing RDS instance

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DAS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `ops_datastore_type` - The database type
* `ops_group_name` - The instance group name
* `ops_group_description` - The description of the instance group
* `ops_group_instance_ids` - The list of instance IDs to be assigned to the group
* `ops_email_template_name` - The name of the email template
* `ops_email_health_rank` - The list of health ranks
* `ops_email_inspection_time` - The diagnosis time
* `ops_email_send_time` - The send time
* `ops_email_time_zone` - The time zone
* `ops_email_subscribe` - Whether to subscribe to the email templates

#### Optional Variables

* `ops_email_address` - The email address for notification (default: null)
* `ops_email_topic` - The topic ID for notification (default: null)
* `ops_email_topic_urn` - The topic URN for notification (default: null)
* `ops_email_obs_bucket_name` - The OBS bucket name for storing inspection reports (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  ops_datastore_type        = "MySQL"
  ops_group_name            = "your_group_name"
  ops_group_description     = "your_group_description"
  ops_group_instance_ids    = ["your_instance_id_1", "your_instance_id_2"]
  ops_email_template_name   = "your_email_template"
  ops_email_health_rank     = ["dangerous", "sub_healthy"]
  ops_email_inspection_time = "00:00-00:00"
  ops_email_send_time       = "08:00-10:00"
  ops_email_time_zone       = "Asia/Shanghai"
  ops_email_subscribe       = true
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
* The `huaweicloud_das_instance_group_assign` is a one-time action resource
* The `huaweicloud_das_email_templates_batch_action` is a one-time action resource
* The email template groups parameter requires instance group IDs from
  `huaweicloud_das_instance_group`
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.92.0 |

# Create a CTS Notification

This example provides best practice code for using Terraform to create a CTS (Cloud Trace Service) notification in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CTS service enabled in the target region
* SMN service enabled in the target region
* IAM agency configured to allow CTS to access SMN resources

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CTS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `topic_name` - The name of the SMN topic for sending notifications
* `notification_name` - The name of the CTS notification
* `notification_agency_name` - The name of the IAM agency which allows CTS to access the SMN resources

#### Optional Variables

* `notification_operation_type` - The type of operation (default: "customized")
* `notification_filter` - The filter conditions for the notification (default: [])
  + `condition` - The logical operator for the filter rules (AND, OR) - **Required when using this list**
  + `rule` - The list of filter rules to apply - **Required when using this list**
* `notification_operations` - The operations to be monitored (default: [])
  + `service` - The name of the service to monitor - **Required when using this list**
  + `resource` - The name of the resource to monitor - **Required when using this list**
  + `trace_names` - The list of trace names to monitor - **Required when using this list**
* `notification_operation_users` - The operation users and groups for the notification (default: [])
  + `group` - The name of the user group - **Required when using this list**
  + `users` - The list of user names in the group - **Required when using this list**

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  topic_name               = "tf_test_topic"
  notification_name        = "tf_test_notification"
  notification_agency_name = "cts_admin_trust"
  
  notification_filter = [
    {
      condition = "OR"
      rule      = [
        "code = 400",
        "resource_name = name",
        "api_version = 1.0",
      ]
    }
  ]
  
  notification_operations = [
    {
      service     = "ECS"
      resource    = "ecs"
      trace_names = [
        "createServer",
        "deleteServer",
      ]
    }
  ]
  
  notification_operation_users = [
    {
      group = "devops"
      users = [
        "your_operation_user_name",
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The CTS notification is dependent on the SMN topic
* The IAM agency must be pre-configured to allow CTS to access SMN resources
* Filter conditions support logical operators (AND, OR) and various rule types
* Operations can be configured to monitor specific services, resources, and trace names
* Operation users can be specified by group or individual user names

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.71.1 |

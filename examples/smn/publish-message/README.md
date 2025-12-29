# Publish SMN Message

This example provides best practice code for using Terraform to publish SMN (Simple Message Notification) messages
on HuaweiCloud. The example demonstrates how to create an SMN topic, add a subscription, and publish messages.

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

* `topic_name` - The name of the SMN topic
* `subscription_protocol` - The protocol of the subscription
* `subscription_endpoint` - The endpoint of the subscription
* `pulblish_subject` - The subject of the message

#### Optional Variables

* `topic_display_name` - The display name of the SMN topic (default: "")
* `enterprise_project_id` - The ID of the enterprise project (default: null)
* `subscription_description` - The description for SMN subscription (default: "")
* `template_name` - The name of the message template (default: "")
* `template_protocol` - The protocol of the message template (default: "", required if template_name is provided)
* `template_content` - The content of the message template (default: "", required if template_name is provided)
* `pulblish_message` - The message content (default: "", mutually exclusive with message_structure and
  message_template_name)
* `pulblish_message_structure` - The JSON message structure that allows sending different content to different protocol
  subscribers (default: "", mutually exclusive with message and message_template_name, required if both template_name and
  pulblish_message are not provided)
* `pulblish_time_to_live` - The maximum retention time of the message within the SMN system in seconds (default: null,
  default value: 3600, max: 86400)
* `pulblish_tags` - The tags of the message (default: {})
* `pulblish_message_attributes` - The message attributes for filtering subscribers (default: [])
  + `name` - The name of the message attribute (required)
  + `type` - The type of the message attribute: STRING, STRING_ARRAY, or PROTOCOL (required)
  + `value` - The value for STRING type (optional)
  + `values` - The values for STRING_ARRAY or PROTOCOL type (optional)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  **Method 1: Publish JSON message structure (send different content to different protocols)**

  ```hcl
  topic_name                 = "tf_test_topic"
  subscription_protocol      = "sms"
  subscription_endpoint      = "your_endpoint"
  pulblish_subject           = "tf_test_subject"
  pulblish_message_structure = "{\"default\":\"Dear user, this is a default message.\",\"sms\":\"Dear user, this is an SMS message.\"}"
  ```

  **Method 2: Publish simple message**

  ```hcl
  topic_name            = "tf_test_topic"
  subscription_protocol = "sms"
  subscription_endpoint = "your_endpoint"
  pulblish_subject      = "tf_test_subject"
  pulblish_message      = "Dear user, this is a simple message."
  ```

  **Method 3: Publish message using template**

  ```hcl
  topic_name            = "tf_test_topic"
  subscription_protocol = "sms"
  subscription_endpoint = "your_endpoint"
  pulblish_subject      = "tf_test_subject"
  template_name         = "tf_test_template"
  template_protocol     = "sms"
  template_content      = "Dear {{name}}, this is a test message."
  pulblish_tags         = {
    name = "John"
  }
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.69.1 |

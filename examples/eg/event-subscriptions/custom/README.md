# Create a custom event subscription

This example provides best practice code for using Terraform to create custom Event Grid (EG) resources in
HuaweiCloud including custom event channels, sources, and event subscriptions.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EG is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `channel_name` - The name of the custom event channel
* `source_name` - The name of the custom event source
* `subscription_name` - The name of the event subscription
* `target_url` - The target url of the event target

#### Optional Variables

* `source_type` - The type of the custom event source (default: "APPLICATION")
* `connection_name` - The exact name of the connection to be queried (default: "default")
* `source_op` - The operation of the source (default: "StringIn")
* `sources_provider_type` - The provider type of the event source (default: "CUSTOM")
* `targets_name` - The name of the event target (default: "HTTPS")
* `targets_provider_type` - The type of the event target (default: "CUSTOM")
* `transform` - The transform configuration of the event target, in JSON format (default: {"type" : "ORIGINAL"})
* `detail_name` - The name(key) of the target detail configuration (default: "detail")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  channel_name      = "your_channel_name"
  source_name       = "your_source_name"
  subscription_name = "your_subscription_name"
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
* The creation of resources may take several minutes depending on the configuration
* This example creates custom EG resources without any dependency on other cloud services
* The example includes a 3-second delay to ensure the custom event channel and source are ready before creating the
  subscription
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.77.4 |

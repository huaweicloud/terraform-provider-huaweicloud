# Create an event subscription with custom event channel

This example provides best practice code for using Terraform to create Event Grid (EG) event subscription in
HuaweiCloud that routes messages to a custom event channel.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EG service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `channel_name` - The name of the custom event channel
* `target_project_id` - The ID of the target project

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "172.16.10.0/24")
* `subnet_gateway` - The gateway IP address of the subnet (default: "172.16.10.1")
* `sources_name` - The name of the event source (default: "HC.VPC")
* `sources_provider_type` - The provider type of the event source (default: "OFFICIAL")
* `source_op` - The operation of the source (default: "StringIn")
* `type_op` - The operation of the type (default: StringIn)
* `targets_name` - The name of the event target (default: "HC.EG")
* `targets_provider_type` - The type of the event target (default: "OFFICIAL")
* `transform` - The transform configuration of the event target, in JSON format (default: "type" : "ORIGINAL")
* `detail_name` - The name(key) of the target detail configuration (default: "eg_detail")
* `agency_name` - The name of the agency (default: "EG_AGENCY")
* `subscription_source_values` - The event types to be subscribed from VPC service
  (default: ["VPC:CloudTrace:ApiCall", "VPC:CloudTrace:ConsoleAction", "VPC:CloudTrace:SystemAction"])
* `target_region_name` - The name of the target region. It will use `region_name` instead if `target_region_name` not
  set (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name          = "your_vpc_name"
  subnet_name       = "your_subnet_name"
  channel_name      = "your_channel_name"
  target_project_id = "your_target_project_id"
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
* This example creates a complete infrastructure including VPC, OBS, and EG resources
* The event subscription automatically filters VPC events and routes them to the custom event channel
* All resources will be created in the specified region
* The example includes a 3-second delay to ensure the custom event channel is ready before creating the subscription

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.68.0 |

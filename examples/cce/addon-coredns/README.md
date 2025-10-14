# Create a CCE Addon for CoreDNS

This example provides best practice code for using Terraform to install and configure the CoreDNS add-on for HuaweiCloud
CCE (Cloud Container Engine).

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster (Make sure that the number of available nodes in the cluster is greater than or equal to 2 and
  all instances are running)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CCE cluster is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `addon_version` - The version of the CoreDNS addon template

#### Optional Variables

* `cluster_id` - The ID of the CCE cluster (default: "")
  If not specified, cluster_name must be provided
* `cluster_name` - The name of the CCE cluster (default: "")
  If not specified, cluster_id must be provided
* `addon_template_name` - The name of the CCE addon template (default: "coredns")
* `project_id` - The ID of the project (default: "")
  If not specified, will be automatically detected from the region

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   # CCE Configuration
   cluster_id    = "your_cce_cluster_id"
   addon_version = "1.30.33" # The version of cluster is v1.32
   ```

3. Initialize Terraform:

   ```bash
   $ terraform init
   ```

4. Review the Terraform plan:

   ```bash
   $ terraform plan
   ```

5. Apply the configuration:

   ```bash
   $ terraform apply
   ```

6. To clean up the resources:

   ```bash
   $ terraform destroy
   ```

## Configuration Details

### CoreDNS Parameters

The example automatically configures CoreDNS with the following key parameters:

* **Basic Configuration (by parameter `basic_json`)**:
  - Automatically preserves all original template basic parameters
  - Maintains default CoreDNS settings for optimal cluster DNS resolution

* **Custom Configuration (by parameter `custom_json`)**:
  - Automatically preserves all original template custom parameters
  - Ensures CoreDNS operates with cluster-specific settings

* **Flavor Configuration (by parameter `flavor_json`)**:
  - Automatically preserves all original template flavor parameters
  - Maintains resource allocation and performance settings

Also you can also customize the parameter values ​for this JSON. This practice only introduces the scenario created
based on the template.

## Note

* Make sure to keep your credentials secure and never commit them to version control
* CoreDNS is a critical component for cluster DNS resolution and should be carefully managed
* Monitor DNS resolution performance after installation to ensure optimal cluster operation
* CoreDNS configuration changes may require cluster restart in some cases

## Version Compatibility

Please refer to the [official Huawei Cloud documentation](https://support.huaweicloud.com/intl/en-us/usermanual-cce/cce_10_0129.html#section7)
for detailed version compatibility information between CCE cluster versions and CoreDNS addon versions.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.31.0 |

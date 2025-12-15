# Create multiple operations for CDN Cache Management

This example provides best practice code for using Terraform to manage CDN (Content Delivery Network) cache in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CDN service enabled in the target region
* At least one CDN domain configured in HuaweiCloud

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CDN cache management operations are performed
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Optional Variables

* `refresh_file_urls` - The list of file URLs that need to be refreshed (up to 1000 URLs)
* `preheat_urls` - The list of URLs that need to be preheated (up to 1000 URLs)
* `zh_url_encode` - Whether to encode Chinese characters in URLs before cache refresh/preheat (default: false)
* `enterprise_project_id` - The ID of the enterprise project to which the resource belongs (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  refresh_file_urls = [
    "https://example.com/index.html",
    "https://example.com/about.html"
  ]
  
  preheat_urls = [
    "https://example.com/index.html",
    "https://example.com/hot-product.html"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* Cache refresh operations typically complete within a few minutes
* Cache preheat operations may take longer depending on the number of URLs
* All URLs must belong to domains that are already configured in CDN
* Chinese URL encoding is useful when URLs contain Chinese characters
* Enterprise project ID is required when using sub-account

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.64.3 |

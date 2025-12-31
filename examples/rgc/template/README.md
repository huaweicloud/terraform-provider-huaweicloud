# Create an RGC template example

This example provides best practice code for using Terraform to create RGC (Resource Governance Center) template
within HuaweiCloud. It demonstrates how to create a predefined template or a customized template.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* RGC service enabled in your HuaweiCloud account

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RGC template is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `template_name` - The name of the template
* `template_type` - The type of the template, only **predefined** and **customized** are supported

#### Optional Variables

* `template_description` - The description of the customized template (default: null)  
  This parameter is valid only when `template_type` is **customized**

* `template_body` - The content of the customized template (default: null)  
   it is a zip-type compressed file that has been encoded using base64, and it is valid
   only when `template_type` is **customized**

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name = "tf_test_template"
  template_type = "predefined"
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

## Example Scenarios

### Scenario 1: Create Predefined Template

```hcl
template_name = "tf_test_template"
template_type = "predefined"
```

### Scenario 2: Create Customized Template

```hcl
template_name        = "my_custom_template"
template_type        = "customized"
template_description = "My custom template description"
template_body        = "base64_encoded_zip_file_content"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.1 |

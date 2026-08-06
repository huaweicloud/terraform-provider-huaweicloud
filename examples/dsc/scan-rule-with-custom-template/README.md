# Create a Custom Scan Rule with Template

This example provides best practice code for using Terraform to create a custom scan rule in HuaweiCloud DSC (Data
Security Center). The example first creates a custom security level, a scan template and a classification under the
template, and then creates a regular expression scan rule that is associated with the template, classification and
security level.

## Prerequisites

* A HuaweiCloud account with DSC permissions
* A purchased DSC instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DSC scan rule is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `security_level_name` - The name of the security level
* `scan_template_name` - The name of the scan template
* `classification_name` - The name of the scan template classification
* `scan_rule_name` - The name of the scan rule

#### Optional Variables

* `security_level_color_number` - The color number of the security level displayed on the console (default: 6)
* `security_level_description` - The description of the security level (default: `""`)
* `scan_template_description` - The description of the scan template (default: "Created_by_terraform_script")
* `scan_rule_match_rate` - The match rate of the scan rule (default: 1)
* `scan_rule_description` - The description of the scan rule (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  security_level_name = "tfscanlevel"
  scan_template_name  = "tfscantemplate"
  classification_name = "tfclassification"
  scan_rule_name      = "tfscanrule"
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
* The DSC instance is not created by this example, please purchase it before applying this example
* The `security_level_name`, `scan_template_name` and `scan_rule_name` must not be the same as existing names in the
  current DSC instance
* Changing the variables in this example updates the existing resources instead of creating new ones
* The scan template is created without built-in rules (`add_built_in_rules = false`), the scan rule created by this
  example is the only rule under the template
* The scan rule is associated with the template through the `templates` block, if you want to delete the scan template,
  please remove the association first

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |

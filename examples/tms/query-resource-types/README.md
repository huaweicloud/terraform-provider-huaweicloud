# Query TMS Resource Types

This example provides best practice code for using Terraform to query the list of service names and resource types
supported by HuaweiCloud Tag Management Service (TMS), with support for various filtering scenarios.

## Prerequisites

* A HuaweiCloud account with TMS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the TMS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Filter Variables

All filter variables are optional. If not specified, all resource types will be returned.

#### Optional Variables

* `exact_service_name` - The exact service name to filter (default: `""`)
* `fuzzy_service_name` - The fuzzy service name pattern to filter (default: `""`)
* `fuzzy_resource_type_name` - The fuzzy resource type name pattern to filter (default: `""`)

## Outputs

The example provides two outputs:

* `regex_matched_service_names` - The list of service names that match the filter criteria
  + Based on fuzzy service name matching
* `regex_matched_resource_types` - The list of resource types that match the filter criteria
  + Based on both service name and resource type name matching

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and configure filter variables as needed:

  ```hcl
  # Example: Filter by fuzzy service name and resource type name
  fuzzy_service_name       = "ccm"
  fuzzy_resource_type_name = "certificate"
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

* View the outputs:

  ```bash
  $ terraform output regex_matched_service_names
  $ terraform output regex_matched_resource_types
  ```

## Example Scenarios

### Scenario 1: Query All Resource Types

Query all resource types supported by TMS without any filters:

```hcl
# Leave all filter variables empty or unset
exact_service_name       = ""
fuzzy_service_name       = ""
fuzzy_resource_type_name = ""
```

This will return all service names and resource types registered in TMS.

### Scenario 2: Filter by Exact Service Name

Query resource types for a specific service using exact matching:

```hcl
exact_service_name = "dli"
```

This will return only resource types that belong to the "dli" service.

### Scenario 3: Filter by Fuzzy Service Name

Query resource types for services matching a pattern:

```hcl
fuzzy_service_name = "ccm"
```

This will return resource types from all services whose names match the pattern "ccm" (e.g., "ccm", "ccm-certificate").

### Scenario 4: Filter by Fuzzy Resource Type Name

Query resource types matching a specific resource type pattern:

```hcl
fuzzy_resource_type_name = "certificate"
```

This will return all resource types whose names match the pattern "certificate" (e.g., "certificate", "certificate_template").

### Scenario 5: Combined Filtering

Combine service name and resource type name filters:

```hcl
fuzzy_service_name       = "ccm"
fuzzy_resource_type_name = "certificate"
```

This will return resource types that:

1. Belong to services matching "ccm" pattern
2. Have resource type names matching "certificate" pattern

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The `exact_service_name` filter is more efficient than `fuzzy_service_name` as it filters at the API level
* Fuzzy matching uses regular expressions, so you can use regex patterns for more complex matching
* The `fuzzy_resource_type_name` filter is applied after service name filtering, so it only searches within already
  filtered results
* Service names in the output are automatically deduplicated
* Resource types include information about whether they are global resources
* The `display_name` field provides localized display names for resource types
* All filters are case-sensitive
* Empty strings (`""`) are treated as "no filter" and will return all results
* The data source queries all resource types from TMS, and filtering is performed locally in Terraform

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.57.0 |

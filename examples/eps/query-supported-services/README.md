# Query EPS Services

This example provides best practice code for using Terraform to query the list of services supported by HuaweiCloud
Enterprise Project Management Service (EPS).

## Prerequisites

* A HuaweiCloud account with EPS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EPS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Filter Variables

All filter variables are optional. If not specified, all services and resource types will be returned.

#### Optional Variables

* `exact_locale` - The display language for service names (default: `""`)
* `fuzzy_service_name` - The fuzzy service name to filter (default: `""`)
* `exact_service_name` - The exact service name to filter (default: `""`)

## Outputs

The example provides two outputs:

* `matched_services` - The list of services that match the filter criteria
* `matched_service_by_only_exact_service_name` - The service that match the specified exact service name

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and configure filter variables as needed:

  ```hcl
  # Example: Filter by fuzzy service name and display language
  fuzzy_service_name = "er"
  exact_locale       = "en-us"
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
  $ terraform output matched_services
  $ terraform output matched_service_by_only_exact_service_name
  ```

## Example Scenarios

### Scenario 1: Query All Supported Services

Query all supported services supported by EPS without any filters:

```hcl
# Leave all filter variables empty or unset
exact_locale       = ""
fuzzy_service_name = ""
exact_service_name = ""
```

This will return all supported services supported by EPS.

### Scenario 2: Filter by Fuzzy Service Name

Query services using fuzzy matching at the API level:

```hcl
fuzzy_service_name = "er"
```

This will query services from EPS using the fuzzy service name filter.

### Scenario 3: Filter by Exact Service Name

Query a specific service using exact matching:

```hcl
exact_service_name = "ServiceStage"
```

This will return the "ServiceStage" service.

### Scenario 4: Combined Filtering

Combine fuzzy service name filter with exact service name filter:

```hcl
fuzzy_service_name = "er"
exact_service_name = "ServiceStage"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* Empty strings (`""`) are treated as "no filter" and will return all results

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.82.3 |

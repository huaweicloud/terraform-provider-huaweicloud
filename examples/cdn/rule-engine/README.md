# Create a CDN Rule Engine for Domain

This example provides best practice code for using Terraform to configure CDN rule engine rules in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* At least one CDN domain configured in HuaweiCloud

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CDN rule engine rule is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `domain_name` - The accelerated domain name to which the rule engine rule belongs
* `rule_name` - The name of the rule engine rule (1-50 characters)

#### Optional Variables

* `rule_status` - Whether to enable the rule engine rule (default: "on")
  Valid values: **on**, **off**
* `rule_priority` - The priority of the rule engine rule (default: 1, range: 1-100)
* `conditions` - The trigger conditions of the rule engine rule, in JSON format (default: "")
* `cache_rule` - The cache rule configuration (default: null)
* `access_control` - The access control configuration (default: null)
* `http_response_headers` - The list of HTTP response header configurations (default: [])
* `browser_cache_rule` - The browser cache rule configuration (default: null)
* `request_url_rewrite` - The access URL rewrite configuration (default: null)
* `flexible_origins` - The list of flexible origin configurations (default: [])
* `origin_request_headers` - The list of origin request header configurations (default: [])
* `origin_request_url_rewrite` - The origin request URL rewrite configuration (default: null)
* `origin_range` - The origin range configuration (default: null)
* `request_limit_rule` - The request rate limit configuration (default: null)
* `error_code_cache` - The list of error code cache configurations (default: [])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  domain_name   = "example.com"
  rule_name     = "api-cache-rule"
  rule_status   = "on"
  rule_priority = 1
  
  # Define conditions
  conditions = <<-JSON
  {
    "match": {
      "logic": "and",
      "criteria": [
        {
          "match_target_type": "path",
          "match_type": "contains",
          "match_pattern": ["/api/"],
          "negate": false,
          "case_sensitive": true
        }
      ]
    }
  }
  JSON

  # Configure cache rule
  cache_rule = {
    ttl           = 10
    ttl_unit      = "m"
    follow_origin = "min_ttl"
  }
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
* Rule engine rules are processed in priority order
* Each action type must be declared in a separate actions block
* Conditions are specified in JSON format and must follow the API specification
* The domain_name cannot be updated after creation (NonUpdatable)
* Rule priority must be unique within the same domain
* All resources will be created in the specified region
* Flexible origin configurations support multiple origin types: ipaddr, domain, obs_bucket, third_bucket
* Error code caching can help reduce origin server load for frequently occurring errors

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.80.5 |

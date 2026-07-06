# Create an AAD instance with black/white list protection

This example provides best practice code for using Terraform to create an Advanced Anti-DDoS (AAD) instance and  
configure black/white lists within HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Understanding of AAD service regions and access types

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the AAD instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `blacklist_ips` - The list of IP addresses to add to the blacklist
* `whitelist_ips` - The list of IP addresses to add to the whitelist
* `instance_config` - The configuration for creating a new AAD instance. Required when `instance_id` is omitted

**instance_config object fields:**

* `instance_name` - The name of the AAD instance
* `service_bandwidth` - The service bandwidth in Mbps
* `ip_type` - The IP type. Valid values:
  + **0** - IPv4
  + **1** - IPv6
* `resource_region` - The resource region. Valid values:
  + **north_china** - North China
  + **east_china** - East China
  + **asia_pacific** - Asia Pacific
* `instance_access_type` - The access type. Valid values:
  + **0** - Website
  + **1** - IP access
* `duration` - The subscription duration
* `amount` - The number of instances to purchase
* `period_type` - The subscription period type. Valid values:
  + **2** - Month
  + **3** - Year

#### Conditionally Required Variables

* `instance_id` - The ID of an existing AAD instance. Required when you already have an AAD instance
* `instance_config` - The configuration for creating a new AAD instance. Required when `instance_id` is omitted

  -> Configure either `instance_id` or `instance_config`. When `instance_id` is provided, the instance resource is
  skipped and `instance_config` is not used.

**instance_config object fields:**

* `basic_bandwidth` - The basic bandwidth in Mbps. Required for **north_china** and **east_china**,
  must be empty for **asia_pacific** (default: null)
* `elastic_bandwidth` - The elastic bandwidth in Mbps. Required for **north_china** and **east_china**,
  must be empty for **asia_pacific** (default: null)
* `basic_qps` - The service QPS. Required when `instance_access_type` is **0**,
  must be empty when `instance_access_type` is **1** (default: null)
* `forwarding_rule` - The number of forwarding rules. Required when `instance_access_type` is **1**,
  must be empty when `instance_access_type` is **0** (default: null)
* `protected_domain` - The number of protected domains. Required when `instance_access_type` is **0**,
  must be empty when `instance_access_type` is **1** (default: null)
* `elastic_service_bandwidth_type` - The elastic service bandwidth type (default: null)  
  Valid values:
  + **2** - Daily 95th percentile
  + **3** - Monthly 95th percentile
* `elastic_service_bandwidth` - The elastic service bandwidth increment (default: null)  
* `protection_package` - The protection package (default: null)  
  Valid values:
  + **basic** - Insurance protection
  + **unlimited** - Unlimited protection
* `enterprise_project_id` - The enterprise project ID to associate (default: null)

## Architecture Overview

This example follows an AAD black/white list configuration workflow with conditional instance creation:

1. **Resolve instance**:

   + When `instance_id` is omitted, create a new AAD instance using `instance_config`
   + When `instance_id` is provided, skip instance creation and use the existing instance

2. **Configure black/white lists**:

   + Create blacklist to block malicious IP addresses
   + Create whitelist to allow trusted IP addresses

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables.

1. Example A — create a new instance with black/white lists:

    ```hcl
    instance_config = {
      instance_name        = "my-aad-instance"
      ip_type              = 0
      resource_region      = "north_china"
      instance_access_type = "1"
      duration             = 1
      amount               = 1
      period_type          = 2
      service_bandwidth    = 100
      basic_bandwidth      = 10
      elastic_bandwidth    = 100
      forwarding_rule      = 50
    }
    
    blacklist_ips = ["192.169.1.100", "192.169.1.101"]
    whitelist_ips = ["10.0.0.1", "10.0.0.2"]
    ```

2. Example B — use an existing instance:

    ```hcl
    instance_id   = "your-existing-aad-instance-id"
    blacklist_ips = ["192.169.1.100"]
    whitelist_ips = ["10.0.0.1"]
    ```

3. Example C — website access type:

    ```hcl
    instance_config = {
      instance_name        = "my-aad-instance-web"
      ip_type              = 0
      resource_region      = "north_china"
      instance_access_type = "0"
      duration             = 1
      amount               = 1
      period_type          = 2
      service_bandwidth    = 100
      basic_bandwidth      = 10
      elastic_bandwidth    = 100
      basic_qps            = 5000
      protected_domain     = 100
    }
    
    blacklist_ips = ["192.169.1.100"]
    whitelist_ips = ["10.0.0.1"]
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

* **IMPORTANT**: AAD instances cannot be deleted or unsubscribed. The `terraform destroy` command will only remove the
  resource from Terraform state, not the actual AAD instance.
* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.
* Configure either `instance_id` or `instance_config`:
  + When `instance_id` is provided, the instance creation is skipped
  + When `instance_config` is provided, a new instance is created
  + When both are provided, using the `instance_id` and will not create a new instance
* The `instance_config.resource_region` parameter determines which configuration parameters are required:
  + **north_china/east_china**: Requires `basic_bandwidth` and `elastic_bandwidth`
  + **asia_pacific**: Requires `protection_package`, cannot have `basic_bandwidth` or `elastic_bandwidth`
* The `instance_config.instance_access_type` parameter determines which configuration parameters are required:
  + **Website ("0")**: Requires `basic_qps` and `protected_domain`
  + **IP access ("1")**: Requires `forwarding_rule`
* The resource uses `lifecycle.ignore_changes` to handle non-updatable parameters that may be missing from
  API response during import.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.9.0  |
| huaweicloud | >= 1.93.0 |

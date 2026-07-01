# Deploy a ModelArts service

This example provides best practice code for using Terraform to deploy a real-time service on ModelArts in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* ModelArts service enabled in the target region
* A dedicated resource pool for inference deployment
* An inference image uploaded to SWR

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ModelArts service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `service_name` - The name of the online inference service
* `service_version` - The version of the online inference service
* `service_group_pool_id` - The ID of the dedicated resource pool for the instance group
* `service_group_name` - The name of the instance group
* `service_unit_configs` - The unit configurations of the instance group
  - `image` - The image configuration of the group unit (required)
    * `source` - The image type
    * `swr_path` - The SWR storage path of the group unit
    * `id` - The image ID (optional)
  - `role` - The role of the group unit (optional)
  - `custom_spec` - The custom resource specification (optional)
    * `memory` - The memory size of the custom specification
    * `cpu` - The CPU number of the custom specification (optional)
    * `gpu` - The GPU number of the custom specification (optional)
    * `ascend` - The number of Ascend chips (optional)
  - `flavor` - The instance flavor of the group unit (optional)
  - `models` - The model mount configuration list (optional, default: [])
    * `source` - The source type
    * `mount_path` - The path to mount into the container
    * `address` - The source address (optional)
    * `source_id` - The source ID (optional)
  - `codes` - The code mount configuration list (optional, default: [])
    * `source` - The source type
    * `mount_path` - The path to mount into the container
    * `address` - The source address (optional)
    * `source_id` - The source ID (optional)
  - `count` - The instance number of the group unit (optional)
  - `cmd` - The startup command of the group unit (optional)
  - `envs` - The environment variables of the group unit (optional)
  - `readiness_health` - The readiness health check configuration (optional)
    * `initial_delay_seconds` - The time to wait when performing the first probe
    * `timeout_seconds` - The timeout for executing the probe
    * `period_seconds` - The period time for performing health check
    * `failure_threshold` - The minimum number of consecutive detection failures
    * `check_method` - The method of the health check
    * `command` - The commands configuration (optional)
    * `url` - The request URL (optional)
  - `startup_health` - The startup health check configuration (optional). Supports the same sub-parameters as
    `readiness_health`
  - `liveness_health` - The liveness health check configuration (optional). Supports the same sub-parameters as
    `readiness_health`
  - `port` - The port of the group unit (optional)
  - `recovery` - The recovery strategy of the group unit (optional)
* `service_runtime_config` - The runtime configuration of the service, in JSON format
  - `service_invoke` - The service invocation configuration
    * `port` - The service port
    * `protocol` - The service protocol
    * `auth_type` - The authentication type
    * `direct_channel_auth_enable` - Whether to enable direct channel authentication
  - `service_limit` - The service limit configuration
    * `request_size_limit` - The request size limit in MB
    * `request_timeout` - The request timeout in seconds
    * `ip_white_list` - The IP whitelist
    * `ip_black_list` - The IP blacklist
    * `rate_limit` - The rate limit configuration
      - `num` - The number of requests
      - `unit` - The time unit
* `service_upgrade_config` - The upgrade configuration of the service, in JSON format
  - `type` - The upgrade type
  - `rolling_update` - The rolling update configuration
    * `max_surge` - The maximum number of pods that can be created above the desired number
    * `max_unavailable` - The maximum number of pods that can be unavailable during the update

#### Optional Variables

* `service_type` - The type of the service (default: "REAL_TIME")
* `service_deploy_type` - The deploy type of the service (default: null)
* `service_description` - The description of the online inference service (default: "")
* `service_group_framework` - The algorithm framework of the instance group (default: "COMMON")
* `service_group_weight` - The weight percentage of the instance group (default: 100)
* `service_group_count` - The number of service instances in the deployment scenario (default: 1)
* `service_log_configs` - The log configurations of the service (default: [])
  - `type` - The type of LTS configuration
  - `log_group_id` - The ID of the LTS group (optional)
  - `log_stream_id` - The ID of the LTS stream (optional)
* `service_tags` - The key/value tags to associate with the service (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  service_name          = "tf_test_service"
  service_version       = "0.0.1"
  service_group_pool_id = "your_resource_pool_id"
  service_group_name    = "tf_test_deploy_group"
  service_unit_configs  = [
    {
      count    = 1
      recovery = "INSTANCE"
      role     = "COMMON"
      flavor   = "your_unit_resource_flavor"

      image = {
        source   = "SWR"
        swr_path = "your-org/your-image:tag"
      }
    }
  ]

  service_runtime_config = <<-JSON
  {
    "service_invoke": {
      "port": 8080,
      "protocol": "HTTPS",
      "auth_type": "TOKEN",
      "direct_channel_auth_enable": false
    },
    "service_limit": {
      "request_size_limit": 20,
      "request_timeout": 30,
      "ip_white_list": [],
      "ip_black_list": [],
      "rate_limit": {
        "num": 200,
        "unit": "SECONDS"
      }
    }
  }
  JSON

  service_upgrade_config = <<-JSON
  {
    "type": "ROLLING",
    "rolling_update": {
      "max_surge": "50%",
      "max_unavailable": "50%"
    }
  }
  JSON
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
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.94.0 |

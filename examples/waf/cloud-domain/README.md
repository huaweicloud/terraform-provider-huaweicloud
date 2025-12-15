# Create a WAF cloud mode domain

This example provides best practice code for using Terraform to create a cloud mode domain in HuaweiCloud
WAF service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the WAF cloud mode domain is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `cloud_instance_resource_spec_code` - The resource specification code of the WAF cloud instance
* `cloud_instance_charging_mode` - The charging mode of the WAF cloud instance
* `cloud_instance_period_unit` - The unit of the subscription period
* `cloud_instance_period` - The subscription period
* `cloud_domain` - The domain name to be protected by WAF
* `cloud_server` - The origin server configurations for the WAF domain
  + `client_protocol` - The protocol type of the client. The options include **HTTP** and **HTTPS**
  + `server_protocol` - The protocol used by WAF to forward client requests to the server. The options include **HTTP**
    and **HTTPS**
  + `address` - The IP address or domain name of the web server that the client accesses
  + `port` - The port number used by the web server. The value ranges from `0` to `65,535`, for example, `8,080`
  + `type` - The type of the origin server. Valid values are: **ipv4** and **ipv6**
  + `weight` - The load balancing algorithm will assign requests to the origin site according to this weight.
    Defaults to `1`.

#### Optional Variables

* `enterprise_project_id` - The enterprise project ID (default: "0")
* `cloud_instance_auto_renew` - Whether to enable auto-renewal for the WAF cloud instance (default: "false")
* `cloud_instance_bandwidth_expack_product` - The configuration of the bandwidth extended packages (default: [])
  + `resource_size` - The size of the bandwidth extended package
* `cloud_instance_domain_expack_product` - The configuration of the domain extended packages (default: [])
  + `resource_size` - The size of the domain extended package
* `cloud_instance_rule_expack_product` - The configuration of the rule extended packages (default: [])
  + `resource_size` - The size of the rule extended package
* `cloud_certificate_id` - The ID of the SSL certificate for the domain (default: "")
* `cloud_certificate_name` - The name of the SSL certificate for the domain (default: "")
* `cloud_proxy` - Whether to enable proxy for the WAF domain (default: false)
* `cloud_description` - The description of the WAF domain (default: "")
* `cloud_website_name` - The website name for the WAF domain (default: "")
* `cloud_protect_status` - The protection status of the WAF domain (default: 0)
* `cloud_forward_header_map` - The field forwarding configuration (default: {})
* `cloud_custom_page` - Configuration for custom error pages (default: [])
  + `http_return_code` - The HTTP return code
  + `block_page_type` - The type of the block page
  + `page_content` - The content of the block page
* `cloud_timeout_settings` - Timeout settings for the WAF domain (default: [])
  + `connection_timeout` - The connection timeout in seconds
  + `read_timeout` - The read timeout in seconds
  + `write_timeout` - The write timeout in seconds
* `cloud_traffic_mark` - Traffic marking configuration for the WAF domain (default: [])
  + `ip_tags` - The IP tags
  + `session_tag` - The session tag
  + `user_tag` - The user tag

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

```hcl
cloud_instance_resource_spec_code = "detection"
cloud_instance_charging_mode      = "prePaid"
cloud_instance_period_unit        = "month"
cloud_instance_period             = 1
cloud_domain                      = "demo-example-test.huawei.com"
cloud_server = [
  {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "119.8.0.17"
    port            = "8080"
    type            = "ipv4"
    weight          = 1
  }
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.61.0 |

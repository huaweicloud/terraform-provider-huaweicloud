# Create a CDN domain

This example provides best practice code for using Terraform to create a CDN (Content Delivery Network) domain in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* SSL certificate files (if using custom certificate for HTTPS)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CDN domain is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `domain_name` - The name of the CDN domain to be accelerated
* `origin_server` - The origin server address (IP address or domain name)

#### Optional Variables

* `domain_type` - The business type of the domain (default: "web")
  Valid values: **web**, **download**, **video**, **wholeSite**
* `service_area` - The area covered by the acceleration service (default: "mainland_china")
  Valid values: **mainland_china**, **outside_mainland_china**, **global**
* `origin_type` - The origin server type (default: "ipaddr")
  Valid values: **ipaddr**, **domain**, **obs_bucket**
* `http_port` - The HTTP port of the origin server (default: 80)
* `https_port` - The HTTPS port of the origin server (default: 443)
* `origin_protocol` - The protocol used to retrieve data from the origin server (default: "http")
  Valid values: **http**, **https**, **follow**
* `ipv6_enable` - Whether to enable IPv6 (default: false)
* `range_based_retrieval_enabled` - Whether to enable range-based retrieval (default: false)
* `domain_description` - The description of the CDN domain (default: "")
* `https_enabled` - Whether to enable HTTPS (default: false)
* `certificate_name` - The name of the SSL certificate (required when https_enabled is true)
* `certificate_source` - The source type of the SSL certificate (required when https_enabled is true)
  Valid values: **0** (HuaweiCloud-managed certificate), **2** (custom certificate)
* `certificate_body_path` - The file path to the SSL certificate content (required when https_enabled is true)
* `private_key_path` - The file path to the private key of the SSL certificate (required when https_enabled is true)
* `http2_enabled` - Whether to enable HTTP/2 (default: false, only valid when https_enabled is true)
* `ocsp_stapling_status` - The OCSP stapling status (default: "off", only valid when https_enabled is true)
  Valid values: **on**, **off**
* `cache_rules` - The cache rules configuration (default: [])
  Each rule object supports the following attributes:
  - `rule_type` - The rule type (required)
    Valid values: **all** (match all files, default), **file_extension** (match by file suffix),
    **catalog** (match by directory), **full_path** (match by full path), **home_page** (match homepage)
  - `content` - The content that matches the rule_type (optional)
    + If `rule_type` is **all** or **home_page**, keep this parameter empty
    + If `rule_type` is **file_extension**, specify file extensions starting with a period (.) and separated by
      semicolons (;), e.g., `.jpg;.zip;.exe` (up to 20 file types)
    + If `rule_type` is **catalog**, specify directories starting with a slash (/) and separated by semicolons (;),
      e.g., `/test/folder01;/test/folder02` (up to 20 directories)
    + If `rule_type` is **full_path**, specify a full path starting with a slash (/) and cannot end with an asterisk,
      e.g., `/test/index.html` or `/test/*.jpg`
  - `ttl` - The cache age in the unit specified by `ttl_type` (required)
    Maximum cache age is 365 days
  - `ttl_type` - The unit of the cache age (required)
    Valid values: **s** (second), **m** (minute), **h** (hour), **d** (day)
  - `priority` - The priority weight of this rule (required)
    Default value is 1. A larger value indicates a higher priority. The value ranges from 1 to 100.
    The weight values must be unique
  - `url_parameter_type` - The URL parameter type (optional, default: "full_url")
    Valid values: **del_params** (ignore specific URL parameters), **reserve_params** (retain specific URL parameters),
    **ignore_url_params** (ignore all URL parameters), **full_url** (retain all URL parameters)
  - `url_parameter_value` - The URL parameter values, separated by commas (,) (optional)
    Up to 10 parameters can be set. Required when `url_parameter_type` is **del_params** or **reserve_params**
* `domain_tags` - The tags of the CDN domain (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  domain_name           = "example.com"
  origin_server         = "192.168.1.100"
  # Enable HTTPS with custom certificate
  https_enabled         = true
  certificate_name      = "terraform-test-cert"
  certificate_source    = "2"
  certificate_body_path = "path/to/your/certificate.crt"
  private_key_path      = "path/to/your/private.key"
  http2_enabled         = true

  # Configure cache rules
  cache_rules = [
    {
      rule_type          = "all"
      content            = ""
      ttl                = 2592000
      ttl_type           = "s"
      priority           = 1
      url_parameter_type = "full_url"
    },
    {
      rule_type          = "file_extension"
      content            = ".jpg;.png;.css;.js"
      ttl                = 604800
      ttl_type           = "s"
      priority           = 2
      url_parameter_type = "ignore_url_params"
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

## Service Area Options

* **mainland_china**: Accelerate content within mainland China
* **outside_mainland_china**: Accelerate content outside mainland China
* **global**: Accelerate content globally

## Domain Type Options

* **web**: Accelerate for the website
* **download**: Accelerate for file downloads
* **video**: Accelerate for on-demand video
* **wholeSite**: Accelerate for the entire site

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The CDN domain creation may take a few minutes to complete
* Before updating the domain configuration, please make sure that the status value is **online**
* The service area cannot be changed between Chinese mainland and outside Chinese mainland
* SSL certificate files should be kept secure and never committed to version control
* Cache rules are processed in priority order (lower number = higher priority)
* All resources will be created in the specified region
* Domain names must be unique within your HuaweiCloud account

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.64.3 |

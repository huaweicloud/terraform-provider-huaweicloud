---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_history_versions"
description: |-
  Use this data source to get the list of API history versions within HuaweiCloud.
---
# huaweicloud_apig_api_history_versions

Use this data source to get the list of API history versions within HuaweiCloud.

## Example Usage

### Query all history versions for the specified API

```hcl
variable "instance_id" {}
variable "api_id" {}

data "huaweicloud_apig_api_history_versions" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
}
```

### Query all history versions for the specified API with environment filter

```hcl
variable "instance_id" {}
variable "api_id" {}
variable "environment_id" {}

data "huaweicloud_apig_api_history_versions" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  env_id      = var.environment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the API history versions are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API belongs.

* `api_id` - (Required, String) Specifies the ID of the API.

* `env_id` - (Optional, String) Specifies the ID of the environment.

* `env_name` - (Optional, String) Specifies the name of the environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_versions` - The list of the API history versions that matched filter parameters.  
  The [api_versions](#apig_api_versions) structure is documented below.

<a name="apig_api_versions"></a>
The `api_versions` block supports:

* `id` - The ID of the API history version.

* `number` - The version number of the API.

* `api_id` - The ID of the API.

* `env_id` - The ID of the published environment.

* `env_name` - The name of the published environment.

* `remark` - The publish description.

* `publish_time` - The publish time of the version, in RFC3339 format.

* `status` - The status of the API version.  
  + **1**: current version is online.
  + **2**: this version is offline.

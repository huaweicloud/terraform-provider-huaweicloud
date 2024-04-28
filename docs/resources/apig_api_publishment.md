---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_publishment"
description: ""
---

# huaweicloud_apig_api_publishment

Using this resource to publish an API to the environment or manage a historical publish version within HuaweiCloud.

~> If you republish on the same environment or switch versions through other ways (such as console) after the API is
published through terraform, the current resource attributes will be affected, resulting in data inconsistency.

## Example Usage

### Publish a new version of the API

```hcl
variable "instance_id" {}
variable "env_id" {}
variable "api_id" {}

resource "huaweicloud_apig_api_publishment" "default" {
  instance_id = var.instance_id
  env_id      = var.env_id
  api_id      = var.api_id
}
```

### Switch to a specified version of the API which is published

```hcl
variable "instance_id" {}
variable "env_id" {}
variable "api_id" {}
variable "version_id" {}

resource "huaweicloud_apig_api_publishment" "default" {
  instance_id = var.instance_id
  env_id      = var.env_id
  api_id      = var.api_id
  version_id  = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to publish APIs.  
  If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the API belongs
  to. Changing this will create a new publishment resource.

* `env_id` - (Required, String, ForceNew) Specifies the ID of the environmentto which the current version of the API
  will be published or has been published.  
  Changing this will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the ID of the API to be published or already published.  
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the current publishment.

* `version_id` - (Optional, String) Specifies the version ID of the current publishment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is constructed from the instance ID, environment ID, and API ID, separated by slashes.

* `env_name` - The name of the environment to which the current version of the API is published.

* `published_at` - Time when the current version was published.

* `publish_id` - The publish ID of the API in current environment.

* `histories` - All publish informations of the API.  
  The [object](#publishment_histories) structure is documented below.

<a name="publishment_histories"></a>
The `histories` block supports:

* `version_id` - The version ID of the API publishment.

* `description` - The version description of the API publishment.

## Import

The publishments can be imported using their related `instance_id`, `env_id` and `api_id`, separated by slashes, e.g.

```shell
$ terraform import huaweicloud_apig_api_publishment.test <instance_id>/<env_id>/<api_id>
```

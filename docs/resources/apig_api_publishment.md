---
subcategory: "API Gateway (Dedicated APIG)"
---

# huaweicloud_apig_api_publishment

API publish Management within HuaweiCloud.

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
  If omitted, the provider-level region will be used. Changing this will create a new publishment resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the API belongs
  to. Changing this will create a new publishment resource.

* `env_id` - (Required, String, ForceNew) Specifies the environment ID to which the current version of the API will be
  published or has been published. Changing this will create a new publishment resource.

* `api_id` - (Required, String, ForceNew) Specifies the API ID to be published or already published.
  Changing this will create a new publishment resource.

* `description` - (Optional, String) Specifies the description of the current publishment.

* `version_id` - (Optional, String) Specifies the version ID of the current publishment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID, which is constructed from the instance ID, environment ID, and API ID, separated by slashes.

* `env_name` - Environment name to which the current version of the API is published.

* `publish_time` - Time when the current version was published.

* `histories` - All publish informations of the API. The structure is documented below.

* `publish_id` - The publish ID of the API in current environment.

The `histories` block supports:

* `version_id` - Version ID of the API publishment.

* `description` - Version description of the API publishment.

## Import

APIs can be imported using their `instance_id`, `env_id` and `api_id`, separated by slashes, e.g.

```
$ terraform import huaweicloud_apig_api_publishment.test
9b0a0a2f97aa43afbf7d852e3ba6a6f9/c5b32727186c4fe6b60408a8a297be09/9a3b3484c08545f9b9b0dcb2de0f5b8a
```

---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance"
description: |-
  Manages a SWR enterprise instance resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_instance

Manages a SWR enterprise instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "spec" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_swr_enterprise_instance" "test" {
  name                  = var.name
  spec                  = var.spec
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the instance.

* `spec` - (Required, String, NonUpdatable) Specifies the specification of the instance. Value can be **swr.ee.professional**.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID .

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `anonymous_access` - (Optional, Bool) Specifies whether to enable anonymous access. Default to **false**.

* `delete_dns` - (Optional, Bool) Specifies whether to delete DNS resources when deleting instance. Default to **false**.

* `delete_obs` - (Optional, Bool) Specifies whether to delete OBS bucket when deleting instance. Default to **false**.

* `obs_encrypt` - (Optional, Bool, NonUpdatable) Specifies whether the OBS bucket is encrypted. Default to **false**.

* `encrypt_type` - (Optional, String, NonUpdatable) Specifies the encrypt type.

* `obs_bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

* `description` - (Optional, String, NonUpdatable) Specifies the description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `charge_mode` - Indicates the charge mode of instance.

* `status` - Indicates the instance status.

* `access_address` - Indicates the access address of instance.

* `user_def_obs` - Indicates whether the user specifies the OBS bucket.

* `version` - Indicates the instance version.

* `vpc_cidr` - Indicates the range of available subnets for the VPC.

* `vpc_name` - Indicates the VPC name.

* `subnet_cidr` - Indicates the range of available subnets for the subnet.

* `subnet_name` - Indicates the subnet name.

* `created_at` - Indicates the creation time.

* `expires_at` - Indicates the expired time.

* `updated_at` - Indicates the last update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes.

## Import

The instance can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `delete_obs`, `delete_dns`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_swr_enterprise_instance" "test" {
    ...

  lifecycle {
    ignore_changes = [
      delete_obs, delete_dns,
    ]
  }
}
```

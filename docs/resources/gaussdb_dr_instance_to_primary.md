---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_instance_to_primary"
description: |-
  Use this resource to promote a GaussDB DR instance to primary within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_instance_to_primary

Use this resource to promote a GaussDB DR instance to primary within HuaweiCloud.

-> This resource is a one-time action resource for promoting a disaster recovery instance to primary. Deleting this
   resource will not revert the operation, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_dr_instance_to_primary" "test" {
  instance_id        = var.instance_id
  disaster_type      = "stream"
  is_support_restore = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to operate the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DR instance ID to be promoted to primary.

* `disaster_type` - (Required, String) Specifies the disaster recovery type.
  The value can be **stream** (streaming disaster recovery).

* `is_support_restore` - (Optional, String) Specifies whether to support disaster recovery failback.
  Only supported for database version V2.0-3.200 and later with streaming disaster recovery type.
  The value can be **true** or **false**. Defaults to **false**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

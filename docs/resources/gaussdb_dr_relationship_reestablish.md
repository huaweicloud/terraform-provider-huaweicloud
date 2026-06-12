---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_relationship_reestablish"
description: |-
  Manages a GaussDB DR relationship re-establish resource within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_relationship_reestablish

Manages a GaussDB DR relationship re-establish resource within HuaweiCloud.

-> This resource is a one-time action resource for re-establishing a disaster recovery relationship. Deleting this
   resource will not revert the operation, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_dr_relationship_reestablish" "test" {
  instance_id   = var.instance_id
  disaster_type = "stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to operate the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID to re-establish the DR relationship.

* `disaster_type` - (Required, String) Specifies the disaster recovery type.
  The value can be **stream** (streaming disaster recovery).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_wdr_snapshot"
description: |-
  Use this resource to generate a WDR snapshot for GaussDB within HuaweiCloud.
---

# huaweicloud_gaussdb_wdr_snapshot

Use this resource to generate a WDR snapshot for GaussDB within HuaweiCloud.

-> This resource is a one-time action resource for generating a WDR snapshot. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_wdr_snapshot" "test" {
  instance_id = var.instance_id
}
```

## Argument

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.

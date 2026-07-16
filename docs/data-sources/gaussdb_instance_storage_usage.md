---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_storage_usage"
description: |-
  Use this data source to query the storage usage of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_storage_usage

Use this data source to query the storage usage of a specified GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_instance_storage_usage" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the storage usage. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `used` - The used storage space of the instance, in GB.

* `total` - The total storage space of the instance, in GB.

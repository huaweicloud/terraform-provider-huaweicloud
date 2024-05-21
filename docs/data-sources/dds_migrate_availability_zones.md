---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_migrate_availability_zones"
description: |-
  Use this data source to query availability zones where DDS instance can migrate.
---

# huaweicloud_dds_migrate_availability_zones

Use this data source to query availability zones where DDS instance can migrate.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_migrate_availability_zones" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DDS instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID in UUID format.

* `names` - The names of availability zone where DDS instance can migrate.

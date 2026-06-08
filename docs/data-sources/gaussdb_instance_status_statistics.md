---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_status_statistics"
description: |-
  Use this data source to query GaussDB instance status statistics within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_status_statistics

Use this data source to query GaussDB instance status statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_instance_status_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances_statistics` - The list of instance status statistics.
  The [instances_statistics](#instances_statistics_struct) structure is documented below.

<a name="instances_statistics_struct"></a>
The `instances_statistics` block supports:

* `status` - The instance status.

* `count` - The number of instances in this status.

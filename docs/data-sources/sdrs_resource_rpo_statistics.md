---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_resource_rpo_statistics"
description: |-
  Use this data source to query SDRS resource RPO statistics within HuaweiCloud.
---

# huaweicloud_sdrs_resource_rpo_statistics

Use this data source to query SDRS resource RPO statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_resource_rpo_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `start_time` - (Optional, String) Specifies the start time using to filter the RPO statistics.
  Format: "yyyy-MM-dd HH:mm:ss.SSS", example: **2019-04-01 12:00:00.000**.

* `end_time` - (Optional, String) Specifies the end time using to filter the RPO statistics.
  Format: "yyyy-MM-dd HH:mm:ss.SSS", example: **2019-04-01 12:00:00.000**.

* `resource_type` - (Optional, String) Specifies the resource type. Valid value is **replication**, which indicates
  querying RPO exceedance trend records for replication pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_rpo_statistics` - The list of RPO exceedance trend records for resources.

  The [resource_rpo_statistics](#resource_rpo_statistics_struct) structure is documented below.

<a name="resource_rpo_statistics_struct"></a>
The `resource_rpo_statistics` block supports:

* `created_at` - The creation time. Format: "yyyy-MM-dd HH:mm:ss.SSS".

* `updated_at` - The update time. Format: "yyyy-MM-dd HH:mm:ss.SSS".

* `id` - The ID of the RPO exceedance trend record for the resource.

* `point_time` - The point-in-time for the RPO exceedance trend record. Format: "yyyy-MM-dd HH:mm".

* `resource_num` - The number of resources with RPO exceedance.

* `resource_type` - The type of resource with RPO exceedance.

---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_subscription_resource"
description: |-
  Use this data source to get the subscription resource of SecMaster.
---

# huaweicloud_secmaster_subscription_resource

Use this data source to get the subscription resource of SecMaster.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_subscription_resource" "test" {
  workspace_id = var.workspace_id
  sku          = "FLOW_DATA_BANDWIDTH"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `sku` - (Optional, String) Specifies the SKU type of the subscription resource.
  The valid values are **FLOW_DATA_BANDWIDTH**, **CSS_CAPACITY**, **PAIMON_CAPACITY**, **OBS_CAPACITY**,
  **JOB_CAPACITY**, **AD_HOC_COUNT**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sku_attribute` - The SKU attribute of the subscription resource.

* `upper_limit` - The upper limit of the resource.

* `unit` - The unit of the resource quota (e.g., GB, count, shard, etc.).

* `step` - The step of the quota.

* `used_amount` - The amount of used resource.

* `unused_amount` - The amount of unused resource.

* `version` - The version number.

* `index_storage_upper_limit` - The upper limit of index storage.

* `index_shards_upper_limit` - The upper limit of index shards.

* `index_shards_unused` - The number of unused index shards.

* `partitions_unused` - The number of unused partitions.

* `partition_upper_limit` - The upper limit of partitions.

* `create_time` - The creation time in milliseconds timestamp.

* `update_time` - The update time in milliseconds timestamp.

---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_upgrade_records"
description: |-
  Use this data source to get the list of DWS cluster upgrade records within HuaweiCloud.
---

# huaweicloud_dws_cluster_upgrade_records

Use this data source to get the list of DWS cluster upgrade records within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_upgrade_records" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster upgrade records are located.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster to which the upgrade records belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of cluster upgrade records that matched filter parameters.  
  The [records](#dws_records_struct) structure is documented below.

<a name="dws_records_struct"></a>
The `records` block supports:

* `id` - The ID of the upgrade record.

* `status` - The status of the upgrade record.

* `record_type` - The type of the upgrade record.

* `from_version` - The source version before the upgrade.

* `to_version` - The target version after the upgrade.

* `start_time` - The start time of the upgrade task, in RFC3339 format.

* `end_time` - The end time of the upgrade task, in RFC3339 format.

* `job_id` - The ID of the upgrade job.

* `failed_reason` - The reason why the upgrade failed.

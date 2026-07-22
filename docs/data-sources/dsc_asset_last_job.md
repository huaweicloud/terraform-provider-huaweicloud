---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_last_job"
description: |-
  Use this data source to get the latest scan job details of a DSC asset within HuaweiCloud.
---

# huaweicloud_dsc_asset_last_job

Use this data source to get the latest scan job details of a DSC asset within HuaweiCloud.

## Example Usage

```hcl
variable "asset_id" {}

data "huaweicloud_dsc_asset_last_job" "test" {
  asset_id = var.asset_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `asset_id` - (Required, String) Specifies the asset ID. The asset ID is used to
  uniquely identify the asset to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `job_id` - The job ID.

* `job_name` - The job name.

* `task_id` - The task ID.

* `type` - The job type. Valid values are **OBS**, **DATABASE**, **BIGDATA**, **GIT**,
  **MRS**, **MRS_HIVE**, **LTS**, **UNKNOWN** and **ALL**.

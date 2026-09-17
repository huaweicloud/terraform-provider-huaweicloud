---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_users"
description: |-
  Use this data source to query the user comparison details of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_compare_users

Use this data source to query the user comparison details of a DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}
variable "compare_job_id" {}

data "huaweicloud_drs_compare_users" "test" {
  job_id         = var.job_id
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the user comparison details.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `compare_job_id` - (Required, String) Specifies the ID of the comparison job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of the user comparison information.

* `user_compare_info` - The user comparison information.

  The [user_compare_info](#user_compare_info_struct) structure is documented below.

<a name="user_compare_info_struct"></a>
The `user_compare_info` block supports:

* `id` - The ID of the user comparison information.

* `src_user_name` - The name of the source database account.

* `tar_user_name` - The name of the target database account.

* `src_privileges` - The privileges of the source database account.

* `tar_privileges` - The privileges of the target database account.

* `is_target_existed` - Whether the target database account exists. The value is **true** if exists, **false** otherwise.

* `compare_result` - The comparison result. The meaning of each value is as follows:
  + **INCONSISTENT**: Inconsistent.
  + **UNABLE_TO_COMPARE**: Unable to compare.
  + **CONSISTENT**: Consistent.
  + **TARGET_SCHEMA_NOT_EXIST**: The target database does not exist.
  + **COMPARE_FAILED**: The comparison failed.
  + **COMPARING**: Comparing.
  + **WAITING_COMPARE**: Waiting for comparison.
  + **UNKNOWN**: Unknown error.

* `created_at` - The creation time.

* `updated_at` - The update time.

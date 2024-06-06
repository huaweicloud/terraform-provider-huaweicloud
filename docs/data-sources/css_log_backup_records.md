---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_log_backup_records"
description: |-
  Use this data source to get the list of CSS cluster log backup records.
---

# huaweicloud_css_log_backup_records

Use this data source to get the list of CSS cluster log backup records.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_log_backup_records" "test" {
    cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CSS cluster.

* `job_id` - (Optional, String) Specifies the ID of the log backup job.

* `type` - (Optional, String) Specifies the type of the log backup job.

* `status` - (Optional, String) Specifies the status of the log backup job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of the CSS cluster log backup records.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The ID of the log backup job.

* `type` - The type of the log backup job.
  + **Manual:** Manual backup.
  + **Auto:** Automatic backup.

* `status` - The status of the log backup job.

* `cluster_id` - The ID of the CSS cluster.

* `log_path` - The storage path of backed up logs in the OBS bucket.

* `create_at` - The creation time.

* `finished_at` - The end time.
  If the creation has not been completed, the end time is empty.

* `failed_msg` - The error information.
  If the task did not fail, the value of this parameter is empty.

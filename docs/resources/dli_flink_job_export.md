---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flink_job_export"
description: |-
  Use this resource to export Flink jobs within HuaweiCloud.
---

# huaweicloud_dli_flink_job_export

Use this resource to export Flink jobs within HuaweiCloud.

-> This resource is a one-time action resource for exporting Flink jobs. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "obs_path" {}
variable "job_ids" {
  type = list(int)
}

resource "huaweicloud_dli_flink_job_export" "test" {
  obs_path = var.obs_path
  job_ids  = var.job_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where to export flink jobs.  
  If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `obs_path` - (Required, String, NonUpdatable) Specifies the OBS save path for the exported job file.  
  e.g. **bucket_name/dir1/dir2**.

* `job_ids` - (Required, List, NonUpdatable) Specifies the set of job IDs to be exported.  
  The type of the job can be **sql job**, **opensource sql job** and **jar job**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

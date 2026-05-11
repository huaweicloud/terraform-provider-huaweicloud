---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinksql_job_savepoint_import"
description: |-
  Use this resource to import a savepoint for a Flink SQL job within HuaweiCloud.
---

# huaweicloud_dli_flinksql_job_savepoint_import

Use this resource to import a savepoint for a Flink SQL job within HuaweiCloud.

-> This resource is a one-time action resource for importing Flink SQL job savepoint. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}
variable "savepoint_path" {}

resource "huaweicloud_dli_flinksql_job_savepoint_import" "test" {
  job_id         = var.job_id
  savepoint_path = var.savepoint_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Flink SQL job savepoint is located.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the Flink SQL job.

* `savepoint_path` - (Required, String, NonUpdatable) Specifies the OBS bucket path of the savepoint.
  You must specify to the parent directory of the metadata file, e.g., **obs://bucket_name/file_name/**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

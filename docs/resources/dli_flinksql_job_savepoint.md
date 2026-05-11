---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinksql_job_savepoint"
description: |-
  Use this resource to trigger a savepoint for a Flink SQL job within HuaweiCloud.
---

# huaweicloud_dli_flinksql_job_savepoint

Use this resource to trigger a savepoint for a Flink SQL job within HuaweiCloud.

-> This resource is a one-time action resource for triggering Flink SQL job savepoint. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}
variable "savepoint_path" {}

resource "huaweicloud_dli_flinksql_job_savepoint" "test" {
  job_id         = var.job_id
  action         = "TRIGGER"
  savepoint_path = var.savepoint_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Flink SQL job savepoint is located.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the Flink SQL job.

* `action` - (Required, String, NonUpdatable) Specifies the operation type of the savepoint. The valid value is
**TRIGGER**.

* `savepoint_path` - (Optional, String, NonUpdatable) Specifies the OBS bucket path of the savepoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

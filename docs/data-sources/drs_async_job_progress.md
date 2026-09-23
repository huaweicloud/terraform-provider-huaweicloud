---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_async_job_progress"
description: |-
  Use this data source to query the execution progress of a DRS async task within HuaweiCloud.
---

# huaweicloud_drs_async_job_progress

Use this data source to query the execution progress of a DRS async task within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "async_job_id" {}

data "huaweicloud_drs_async_job_progress" "test" {
  async_job_id = var.async_job_id
}
```

### Filter by Type

```hcl
variable "async_job_id" {}

data "huaweicloud_drs_async_job_progress" "test" {
  async_job_id = var.async_job_id
  type         = "file"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the async job progress.  
  If omitted, the provider-level region will be used.

* `async_job_id` - (Required, String) Specifies the unique identifier of the async task.  
  This ID is obtained from the response of the export or template creation API.

* `type` - (Optional, String) Specifies the type of the async task to query.  
  Defaults to **file**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The status of the async task.

* `progress` - The execution progress of the async task.

* `error_code` - The error code of the async task.

* `error_message` - The error message of the async task.

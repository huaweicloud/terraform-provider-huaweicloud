---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_restart"
description: |-
  Use this resource to restart compare tasks of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_compare_restart

Use this resource to restart compare tasks of a DRS job within HuaweiCloud.

-> This resource is a one-time action resource for restarting compare tasks of a DRS job. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}
variable "compare_job_id" {}

resource "huaweicloud_drs_compare_restart" "test" {
  job_id          = var.job_id
  compare_job_ids = [var.compare_job_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to restart the compare tasks.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job to restart compare tasks.

* `compare_job_ids` - (Required, List, NonUpdatable) Specifies the list of compare job IDs to restart.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the request ID returned by the API.

---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscription_delete"
description: |-
  Manages a resource to delete a DRS subscription job within HuaweiCloud.
---

# huaweicloud_drs_subscription_delete

Manages a resource to delete a DRS subscription job within HuaweiCloud.

-> This resource is a one-time action resource used to delete a DRS subscription job. Deleting this resource will not
   restore the deleted subscription job or undo the delete action, but will only remove the resource information from
   the tf state file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_subscription_delete" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to delete the subscription job.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS subscription job to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_change_flavor"
description: |-
  Use this resource to change the flavor of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_job_change_flavor

Use this resource to change the flavor of a DRS job within HuaweiCloud.

-> This resource is a one-time action resource for changing the flavor of a DRS job. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_job_change_flavor" "test" {
  job_id    = var.job_id
  spec_type = "high"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to change the DRS job flavor.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job to change flavor.

* `spec_type` - (Required, String, NonUpdatable) Specifies the target specification type of the DRS job.  
  The valid values are as follows:
  + **micro**
  + **small**
  + **medium**
  + **high**
  + **xlarge**
  + **2xlarge**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the flavor change task ID.

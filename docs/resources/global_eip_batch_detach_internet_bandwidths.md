---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_batch_detach_internet_bandwidths"
description: ""
---

# huaweicloud_global_batch_detach_internet_bandwidths

Manages a resource to batch detach internet bandwidths from global EIPs within HuaweiCloud.

-> **NOTE:** This resource is a one-time action resource used to batch detach internet bandwidths from global EIPs.
Deleting this resource will not clear the corresponding request record, but will only remove resource information from
the tfstate file.

## Example Usage

```hcl
variable "global_eip_id_1" {}
variable "global_eip_id_2" {}

resource "huaweicloud_global_batch_detach_internet_bandwidths" "test" {
  global_eips {
    global_eip_id         = var.global_eip_id_1
  }

  global_eips {
    global_eip_id         = var.global_eip_id_2
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `global_eips` - (Required, List) Specifies the list of global EIPs to be detached from internet bandwidths.
  The [global_eips](#global_eips_struct) structure is documented below.

<a name="global_eips_struct"></a>
The `global_eips` block supports:

* `global_eip_id` - (Required, String) Specifies the ID of the global EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `job_id` - The ID of the batch detach job.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_batch_attach_internet_bandwidths"
description: |-
  Manages a resource to batch attach internet bandwidths to global EIPs within HuaweiCloud.
---

# huaweicloud_global_batch_attach_internet_bandwidths

Manages a resource to batch attach internet bandwidths to global EIPs within HuaweiCloud.

-> This resource is a one-time action resource used to batch attach internet bandwidths to global EIPs.
  Deleting this resource will not clear the corresponding request record, but will only remove resource information from
  the tfstate file.

## Example Usage

```hcl
variable "global_eip_configurations" {
  type = list(object({
    global_eip_id         = string
    internet_bandwidth_id = string
  }))
}

resource "huaweicloud_global_batch_attach_internet_bandwidths" "test" {
  dynamic "global_eips" {
    for_each = var.global_eip_configurations
    content {
      global_eip_id         = global_eips.value.global_eip_id
      internet_bandwidth_id = global_eips.value.internet_bandwidth_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `global_eips` - (Required, List,NonUpdatable) Specifies the list of global EIPs and internet bandwidths to be attached
  The [global_eips](#global_eips_struct) structure is documented below.

<a name="global_eips_struct"></a>
The `global_eips` block supports:

* `global_eip_id` - (Required, String) Specifies the ID of the global EIP.

* `internet_bandwidth_id` - (Required, String) Specifies the ID of the internet bandwidth.
  All global EIPs in the list must be bound to the same bandwidth ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

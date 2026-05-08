---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_bandwidth_batch_attach"
description: |-
  Manages a resource to batch attach EIPs to shared bandwidth within HuaweiCloud.
---

# huaweicloud_eip_bandwidth_batch_attach

Manages a resource to batch attach EIPs to shared bandwidth within HuaweiCloud.

-> This resource is a one-time action resource for batch attaching EIPs to shared bandwidth.
Deleting this resource will not detach the EIPs from the shared bandwidth, but will only remove
the resource information from the Terraform state file.

## Example Usage

```hcl
variable "eip_bandwidth_attachments" {
  type = list(object({
    bandwidth_id = string
    publicip_id  = string
  }))
}

resource "huaweicloud_eip_bandwidth_batch_attach" "test" {
  dynamic "publicips" {
    for_each = var.eip_bandwidth_attachments
    content {
      bandwidth_id = publicips.value.bandwidth_id
      publicip_id  = publicips.value.publicip_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC EIP associate resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `publicips` - (Required, List, NonUpdatable) Specifies the list of EIPs to attach to shared bandwidth.

The [publicips](#publicips_struct) structure is documented below.

<a name="publicips_struct"></a>
The `publicips` block supports:

* `bandwidth_id` - (Required, String, NonUpdatable) Specifies the shared bandwidth ID.

* `publicip_id` - (Required, String, NonUpdatable) Specifies the EIP ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

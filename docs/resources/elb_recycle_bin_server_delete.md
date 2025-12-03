---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_recycle_bin_loadbalancer_delete"
description: |-
  Manages an ELB recycle bin load balancer delete resource within HuaweiCloud.
---

# huaweicloud_elb_recycle_bin_loadbalancer_delete

Manages an ELB recycle bin load balancer delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}

resource "huaweicloud_elb_recycle_bin_loadbalancer_delete" "test" {
  loadbalancer_id = var.loadbalancer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `loadbalancer_id` - (Required, String, NonUpdatable) Specifies the load balancer ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the load balancer ID.

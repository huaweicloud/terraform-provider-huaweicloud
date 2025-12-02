---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_recycle_bin"
description: |-
  Manages an ELB recycle bin resource within HuaweiCloud.
---

# huaweicloud_elb_recycle_bin

Manages an ELB recycle bin resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_elb_recycle_bin" "test" {
  retention_hour         = 10
  recycle_threshold_hour = 50
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `retention_hour` - (Optional, String) Specifies how long you want your load balancers to be kept in the recycle bin.
  An integer number is required.

* `recycle_threshold_hour` - (Optional, String) Specifies how old a deleted or unsubscribed load balancer has to be for
  it to be moved to the recycle bin. A load balancer can only be moved to the recycle bin if it has been there longer than
  the specified number of days in the policy. An integer number is required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the project ID.

## Import

The ELB recycle bin resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_recycle_bin.test <id>
```

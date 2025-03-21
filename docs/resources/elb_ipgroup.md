---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_ipgroup"
description: |-
  Manages a Dedicated ELB Ip Group resource within HuaweiCloud.
---

# huaweicloud_elb_ipgroup

Manages a Dedicated ELB Ip Group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_elb_ipgroup" "basic" {
  name        = "basic"
  description = "basic example"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ip group resource. If omitted, the
  provider-level region will be used. Changing this creates a new ip group.

* `name` - (Required, String) Human-readable name for the ip group.

* `description` - (Optional, String) Human-readable description for the ip group.

* `ip_list` - (Required, List) Specifies an array of one or more ip addresses. The ip_list object structure is
  documented below. The [ip_list](#ELB_iplistResp) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the ip group. Changing this
  creates a new ip group.

<a name="ELB_iplistResp"></a>
The `ip_list` block supports:

* `ip` - (Required, String) IP address or CIDR block.

* `description` - (Optional, String) Human-readable description for the ip.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The uuid of the ip group.

* `listener_ids` - The listener IDs which the ip group associated with.

* `created_at` - The creation time of the ip group.

* `updated_at` - The update time of the ip group.

## Import

The ELB IP group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_ipgroup.test <id>
```

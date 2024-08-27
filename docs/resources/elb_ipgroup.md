---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_ipgroup"
description: ""
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
  documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the ip group. Changing this
  creates a new ip group.

The `ip_list` block supports:

* `ip` - (Required, String) IP address or CIDR block.

* `description` - (Optional, String) Human-readable description for the ip.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The uuid of the ip group.

* `listener_ids` - The listener IDs which the ip group associated with.

* `created_at` - The create time of the ip group.

* `updated_at` - The update time of the ip group.

ELB IP group can be imported using the IP group ID, e.g.

```bash
$ terraform import huaweicloud_elb_ipgroup.group_1 5c20fdad-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a IP group.
You can then decide if changes should be applied to the IP group, or the resource
definition should be updated to align with the IP group. Also you can ignore changes as below.

```hcl
resource "huaweicloud_elb_ipgroup" "group_1" {
    ...
  lifecycle {
    ignore_changes = [
      enterprise_project_id,
    ]
  }
}
```

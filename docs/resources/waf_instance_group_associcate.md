---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_instance_group_associate

Associate ELB instances to a WAF instance group.

## Example Usage

```hcl
variable "group_id" {}
variable "elb_instance_id" {}

resource "huaweicloud_waf_instance_group_associate" "group_associate" {
  group_id      = group_id
  load_balances = [elb_instance_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which the WAF instance group created.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of instance group.
  Changing this will create a new instance.

* `load_balances` - (Required, List) Specifies the IDs of the ELB instances bound to the instance group.
  This is an array of security group ids.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID in UUID format.

## Import

The instance group associate can be imported using the group ID, e.g.:

```sh
terraform import huaweicloud_waf_instance_group_associate.group_associate 0be1e69d-1987-4d9c-9dc5-fc7eed592398
```

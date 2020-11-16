---
subcategory: "Virtual Private Network (VPN)"
---

# huaweicloud\_vpnaas\_endpoint\_group\_v2

Manages a V2 Endpoint Group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_vpnaas_endpoint_group_v2" "group_1" {
  name = "Group 1"
  type = "cidr"
  endpoints = ["10.2.0.0/24",
        "10.3.0.0/24",]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an endpoint group. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    group.

* `name` - (Optional) The name of the group. Changing this updates the name of
    the existing group.

* `tenant_id` - (Optional) The owner of the group. Required if admin wants to
    create an endpoint group for another project. Changing this creates a new group.

* `description` - (Optional) The human-readable description for the group.
    Changing this updates the description of the existing group.

* `type` -  (Optional) The type of the endpoints in the group. A valid value is subnet, cidr, network, router, or vlan.
    Changing this creates a new group.

* `endpoints` - (Optional) List of endpoints of the same type, for the endpoint group. The values will depend on the type.
    Changing this creates a new group.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `type` - See Argument Reference above.
* `endpoints` - See Argument Reference above.
* `value_specs` - See Argument Reference above.


## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpnaas_endpoint_group_v2.group_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```

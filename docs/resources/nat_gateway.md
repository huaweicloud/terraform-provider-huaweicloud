---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud\_nat\_gateway

Manages a Nat gateway resource within HuaweiCloud Nat
This is an alternative to `huaweicloud_nat_gateway_v2`

## Example Usage

```hcl
resource "huaweicloud_nat_gateway" "nat_1" {
  name                = "Terraform"
  description         = "test for terraform"
  spec                = "3"
  router_id           = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
  internal_network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Nat gateway resource. If omitted, the provider-level region will work as default. Changing this creates a new Nat gateway resource.

* `name` - (Required) The name of the nat gateway.

* `description` - (Optional) The description of the nat gateway.

* `spec` - (Required) The specification of the nat gateway, valid values are "1",
    "2", "3", "4".

* `tenant_id` - (Optional) The target tenant ID in which to allocate the nat
    gateway. Changing this creates a new nat gateway.

* `router_id` - (Required) ID of the router this nat gateway belongs to. Changing
    this creates a new nat gateway.

* `internal_network_id` - (Optional) ID of the network this nat gateway connects to.
    Changing this creates a new nat gateway.

* `enterprise_project_id` - (Optional) The enterprise project id of the nat gateway. 
    Changing this creates a new nat gateway.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `spec` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `internal_network_id` - See Argument Reference above.

## Import

Nat gateway can be imported using the following format:

```
$ terraform import huaweicloud_nat_gateway.nat_1 d126fb87-43ce-4867-a2ff-cf34af3765d9
```
---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud\_nat\_dnat\_rule

Manages a Dnat rule resource within HuaweiCloud Nat.
This is an alternative to `huaweicloud_nat_dnat_rule_v2`

## Example Usage

### Dnat

```hcl
resource "huaweicloud_nat_dnat_rule" "dnat_1" {
  nat_gateway_id        = "bf99c679-9f41-4dac-8513-9c9228e713e1"
  floating_ip_id        = "2bd659ab-bbf7-43d7-928b-9ee6a10de3ef"
  private_ip            = "10.0.0.12"
  protocol              = "tcp"
  internal_service_port = 993
  external_service_port = 242
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the dnat rule resource.
  If omitted, the provider-level region will be used. Changing this creates a new dnat rule.

* `nat_gateway_id` - (Required, String, ForceNew) ID of the nat gateway this dnat rule belongs to.
   Changing this creates a new dnat rule.

* `floating_ip_id` - (Required, String, ForceNew) Specifies the ID of the floating IP address.
  Changing this creates a new dnat rule.

* `protocol` - (Required, String, ForceNew) Specifies the protocol type. Currently,
  TCP, UDP, and ANY are supported.
  Changing this creates a new dnat rule.

* `internal_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or BMSs
  to provide services for external systems. Changing this creates a new dnat rule.

* `external_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or
  BMSs to provide services for external systems.
  Changing this creates a new dnat rule.

* `port_id` - (Optional, String, ForceNew) Specifies the port ID of an ECS or a BMS.
  This parameter and private_ip are alternative. Changing this creates a
  new dnat rule.

* `private_ip` - (Optional, String, ForceNew) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and port_id are alternative.
  Changing this creates a new dnat rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `created_at` - Dnat rule creation time.

* `status` - Dnat rule status.

* `floating_ip_address` - The actual floating IP address.

## Import

Dnat can be imported using the following format:

```
$ terraform import huaweicloud_nat_dnat_rule.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```

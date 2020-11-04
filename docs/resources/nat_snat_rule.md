---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud\_nat\_snat\_rule

Manages a Snat rule resource within HuaweiCloud Nat
This is an alternative to `huaweicloud_nat_snat_rule_v2`

## Example Usage

```hcl
resource "huaweicloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = "3c0dffda-7c76-452b-9dcc-5bce7ae56b17"
  network_id     = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
  floating_ip_id = "0a166fc5-a904-42fb-b1ef-cf18afeeddca"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the snat rule resource. If omitted, the provider-level region will work as default. Changing this creates a new snat rule resource.

* `nat_gateway_id` - (Required) ID of the nat gateway this snat rule belongs to.
    Changing this creates a new snat rule.

* `network_id` - (Required) ID of the network this snat rule connects to.
    Changing this creates a new snat rule.

* `floating_ip_id` - (Required) ID of the floating ip this snat rule connets to.
    Changing this creates a new snat rule.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `nat_gateway_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `floating_ip_id` - See Argument Reference above.
* `floating_ip_address` - The actual floating IP address.
* `status` - The status of the snat rule.

## Import

Snat can be imported using the following format:

```
$ terraform import huaweicloud_nat_snat_rule.snat_1 9e0713cb-0a2f-484e-8c7d-daecbb61dbe4
```
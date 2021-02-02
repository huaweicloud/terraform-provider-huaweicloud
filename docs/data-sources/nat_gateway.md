---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud\_nat\_gateway

Use this data source to get the information of an available HuaweiCloud NAT gateway.


## Example Usage

```hcl
data "huaweicloud_nat_gateway" "natgateway" {
  name = "tf_test_natgateway"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the Nat
    gateway resource. If omitted, the provider-level region will be used.

* `id` - (Optional, String) Specifies the ID of the NAT gateway.

* `name` - (Optional, String) Specifies the nat gateway name. The name can
    contain only digits, letters, underscores (_), and hyphens(-).

* `internal_network_id` - (Optional, String) Specifies the network ID of the
    downstream interface (the next hop of the DVR) of the NAT gateway.

* `router_id` - (Optional, String) Specifies the ID of the router this nat
    gateway belongs to.

* `spec` - (Optional, String) The NAT gateway type.
    The value can be:
    * `1`: small type, which supports up to 10,000 SNAT connections.
    * `2`: medium type, which supports up to 50,000 SNAT connections.
    * `3`: large type, which supports up to 200,000 SNAT connections.
    * `4`: extra-large type, which supports up to 1,000,000 SNAT connections.

* `description` - (Optional, String) Specifies the description of the nat
   gateway. The value contains 0 to 255 characters, and angle brackets (<)
   and (>) are not allowed.

* `status` - (Optional, String) Specifies the status of the NAT gateway.

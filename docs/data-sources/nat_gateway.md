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

* `region` - (Optional, String) The region in which to obtain the gateways. If omitted, the provider-level region will be used.

* `id` - (Optional, String) The ID of the NAT gateway.

* `name` - (Optional, String) The name of the NAT gateway.

* `description` - (Optional, String) The information about the NAT gateway..

* `spec` - (Optional, String) The NAT gateway type.
              The value can be:
              1: small type, which supports up to 10,000 SNAT connections.
              2: medium type, which supports up to 50,000 SNAT connections.
              3: large type, which supports up to 200,000 SNAT connections.
              4: extra-large type, which supports up to 1,000,000 SNAT connections.

* `router_id` - (Optional, String) The router ID.

* `internal_network_id` - (Optional, String) The network ID of the downstream interface (the next hop of the DVR) of the NAT gateway.

* `status` - (Optional, String) The status of the NAT gateway.


## Attributes Reference

The following attributes
are exported:

* `admin_state_up` - The unfrozen or frozen state.
                        The value can be:
                          true: indicates the unfrozen state.
                          false: indicates the frozen state.

---
subcategory: "Elastic IP (EIP)"
---

# huaweicloud\_vpc\_bandwidth

Provides details about a specific bandwidth.

## Example Usage

```hcl
variable "bandwidth_name" {}

data "huaweicloud_vpc_bandwidth" "bandwidth_1" {
  name = var.bandwidth_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
bandwidth in the current tenant. The following arguments are supported:

* `region` - (Optional) The region in which to obtain the bandwidth. If omitted, the provider-level region will be used.

* `name` - (Required) The name of the Shared Bandwidth to retrieve.

* `size` - (Optional) The size of the Shared Bandwidth to retrieve. The value ranges from 5 to 2000 G.

* `enterprise_project_id` - (Optional) The enterprise project id of the Shared Bandwidth to retrieve.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the Shared Bandwidth.

* `name` -  See Argument Reference above.

* `size` - See Argument Reference above.

* `enterprise_project_id` - See Argument Reference above.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `bandwidth_type` - Indicates the bandwidth type.

* `charge_mode` - Indicates whether the billing is based on traffic, bandwidth, or 95th percentile bandwidth (enhanced).

* `status` - Indicates the bandwidth status.

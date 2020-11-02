---
subcategory: "Elastic IP (EIP)"
---

# huaweicloud\_vpc\_bandwidth

Manages a Shared Bandwidth resource within HuaweiCloud.
This is an alternative to `huaweicloud_vpc_bandwidth_v2`

## Example Usage

```hcl
resource "huaweicloud_vpc_bandwidth" "bandwidth_1" {
  name = "bandwidth_1"
  size = 5
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Shared Bandwidth. If omitted, the provider-level region will work as default. Changing this creates a new Shared Bandwidth resource.

* `name` - (Required) The name of the Shared Bandwidth.

* `size` - (Required) The size of the Shared Bandwidth. The value ranges from 5 to 2000 G.

* `enterprise_project_id` - (Optional) The enterprise project id of the Shared Bandwidth. Changing this creates a new bandwidth.


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

## Import

Shared Bandwidths can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc_bandwidth.bandwidth_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
